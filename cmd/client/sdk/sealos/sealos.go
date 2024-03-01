package sealos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/sdk"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/tools"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type SDK struct {
}

var (
	num    int32 = 1
	prefix       = networkingv1.PathTypePrefix
	fl           = false
)

type Amount struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ActivityBonus           int    `json:"activityBonus"`
		Balance                 int    `json:"balance"`
		DeductionBalance        int    `json:"deductionBalance"`
		EncryptBalance          string `json:"encryptBalance"`
		EncryptDeductionBalance string `json:"encryptDeductionBalance"`
	} `json:"data"`
}

func (s *SDK) Auth(providerId uint) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	amountUrl := fmt.Sprintf("https://costcenter.%s/api/account/getAmount", sdk.SealosRegionMap[*provider.Region])

	req, err := http.NewRequest("GET", amountUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", url.PathEscape(provider.CloudAuth.KubeConfig))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("error request : " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var sa = Amount{}
	err = json.Unmarshal(body, &sa)
	if err != nil {
		return err
	}

	*provider.Amount = float64(sa.Data.Balance-sa.Data.DeductionBalance) / 1000000
	*provider.Cost = float64(sa.Data.DeductionBalance) / 1000000

	service.GetService("provider").Update(provider.ID, provider)

	return nil
}

func (s *SDK) Billing(providerId uint, tunnel models.Tunnel) error {
	// 详细计算某个隧道花费数据
	//url := fmt.Sprintf("https://costcenter.%s/api/billing", SealosRegionMap[provider.Region])
	return nil
}

func (s *SDK) Deploy(providerId uint, tun *models.Tunnel) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	// sealos 部署十分简单，直接调用 k8s client-go 即可。
	ctx := context.TODO()

	ns, clientSet, err := parseKubeConfig(provider.CloudAuth.KubeConfig)

	if err != nil {
		return err
	}

	svcName := "seamoon-" + *tun.Name + "-" + string(*tun.Type)
	imgName := "dvkunion/seamoon:" + consts.Version
	hostName := tools.GenerateRandomLetterString(12)
	port, err := strconv.Atoi(*tun.Port)
	if err != nil {
		return err
	}
	// 需要先创建 deployment 负载
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcName,
			Annotations: map[string]string{
				"originImageName":                    imgName,
				"deploy.cloud.sealos.io/minReplicas": "1",
				"deploy.cloud.sealos.io/maxReplicas": "1",
				"deploy.cloud.sealos.io/resize":      "0Gi",
			},
			Labels: map[string]string{
				"cloud.sealos.io/app-deploy-manager": svcName,
				"app":                                svcName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:             &num,
			RevisionHistoryLimit: &num,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": svcName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": svcName,
					},
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: &fl,
					Containers: []corev1.Container{
						{
							Name:  svcName,
							Image: imgName,
							Env: func() []corev1.EnvVar {
								var env = make([]corev1.EnvVar, 0)
								if tun.TunnelConfig.Tor {
									// 需要增加环境变量
									env = append(env, corev1.EnvVar{
										Name:  "SEAMOON_TOR",
										Value: "true",
									})
								}
								return env
							}(),
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU: func() resource.Quantity {
										if tun.TunnelConfig.CPU < 0.1 {
											return resource.MustParse("10m")
										}
										return resource.MustParse(strconv.Itoa(int(tun.TunnelConfig.CPU*100)) + "m")
									}(),
									corev1.ResourceMemory: func() resource.Quantity {
										if tun.TunnelConfig.Memory < 64 {
											return resource.MustParse("6Mi")
										}
										return resource.MustParse(strconv.Itoa(int(tun.TunnelConfig.Memory/10)) + "Mi")
									}(),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU: func() resource.Quantity {
										if tun.TunnelConfig.CPU < 0.1 {
											return resource.MustParse("100m")
										}
										return resource.MustParse(strconv.Itoa(int(tun.TunnelConfig.CPU*1000)) + "m")
									}(),
									corev1.ResourceMemory: func() resource.Quantity {
										if tun.TunnelConfig.Memory < 64 {
											return resource.MustParse("64Mi")
										}
										return resource.MustParse(strconv.Itoa(int(tun.TunnelConfig.Memory)) + "Mi")
									}(),
								},
							},
							Command: []string{"/app/seamoon"},
							Args: func() []string {
								switch *tun.Type {
								case tunnel.WST:
									return []string{"server", "-p", "9000", "-t", "websocket"}
								case tunnel.GRT:
									return []string{"server", "-p", "8089", "-t", "grpc"}
								}
								return []string{}
							}(),
							Ports: []corev1.ContainerPort{
								{
									Name:          "seamoon-http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: int32(port),
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}

	// 使用客户端创建Deployment
	_, err = clientSet.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Deployment创建成功！")
	// 然后是创建 service 和 ingress
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcName,
			Labels: map[string]string{
				"cloud.sealos.io/app-deploy-manager": svcName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": svcName,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt32(int32(port)),
				},
			},
		},
	}
	_, err = clientSet.CoreV1().Services(ns).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Service创建成功！")
	// ingress
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: svcName,
			Labels: map[string]string{
				"cloud.sealos.io/app-deploy-manager":        svcName,
				"cloud.sealos.io/app-deploy-manager-domain": hostName,
			},
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":                 "nginx",
				"nginx.ingress.kubernetes.io/proxy-body-size": "32m",
				"nginx.ingress.kubernetes.io/ssl-redirect":    "false",
				"nginx.ingress.kubernetes.io/backend-protocol": func() string {
					switch *tun.Type {
					case tunnel.WST:
						return "WS"
					case tunnel.GRT:
						return "GRPC"
					}
					return "HTTP"
				}(),
				"nginx.ingress.kubernetes.io/proxy-send-timeout": "3600",
				"nginx.ingress.kubernetes.io/proxy-read-timeout": "3600",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: fmt.Sprintf("%s.%s", hostName, sdk.SealosRegionMap[*provider.Region]),
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: svcName,
											Port: networkingv1.ServiceBackendPort{
												Number: 9000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{fmt.Sprintf("%s.%s", hostName, sdk.SealosRegionMap[*provider.Region])},
					SecretName: "wildcard-cloud-sealos-io-cert",
				},
			},
		},
	}
	_, err = clientSet.NetworkingV1().Ingresses(ns).Create(ctx, ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Ingress创建成功！")
	*tun.Status = tunnel.ACTIVE
	*tun.Addr = fmt.Sprintf("%s.%s", hostName, sdk.SealosRegionMap[*provider.Region])
	service.GetService("tunnel").Update(tun.ID, tun)
	return nil
}

func (s *SDK) Destroy(providerId uint, tun *models.Tunnel) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	ctx := context.TODO()

	ns, clientSet, err := parseKubeConfig(provider.CloudAuth.KubeConfig)

	if err != nil {
		return err
	}

	svcName := "seamoon-" + *tun.Name + "-" + string(*tun.Type)
	if err := clientSet.AppsV1().Deployments(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	slog.Info("成功删除")
	if err := clientSet.CoreV1().Services(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	slog.Info("成功删除")
	if err := clientSet.NetworkingV1().Ingresses(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	slog.Info("成功删除")

	// 删除数据
	return nil
}

func (s *SDK) SyncFC(providerId uint) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	tunList := make([]models.Tunnel, 0)
	ctx := context.TODO()

	ns, clientSet, err := parseKubeConfig(provider.CloudAuth.KubeConfig)

	if err != nil {
		return err
	}

	svcs, err := clientSet.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	ingresses, err := clientSet.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, svc := range svcs.Items {
		if strings.HasPrefix(svc.Name, "seamoon-") {
			// 说明是我们的服务，继续获取对应的 label 来查找 ingress
			var tun = models.Tunnel{}
			*tun.Name = strings.Split(svc.Name, "-")[1]
			*tun.Port = strconv.Itoa(int(svc.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort))
			*tun.Type = func() tunnel.Type {
				if strings.HasSuffix(svc.Name, "websocket") {
					return tunnel.WST
				}
				if strings.HasSuffix(svc.Name, "grpc") {
					return tunnel.GRT
				}
				return tunnel.NULL
			}()
			for _, ingress := range ingresses.Items {
				if ingress.Name == svc.Name {
					*tun.Addr = ingress.Spec.Rules[0].Host
				}
			}
			t := service.GetService("tunnel").Create(&tun).(*models.Tunnel)
			tunList = append(tunList, *t)
		}
	}
	if len(tunList) > 0 {
		provider.Tunnels = append(provider.Tunnels, tunList...)
	}

	service.GetService("provider").Update(provider.ID, provider)
	return nil
}

func parseKubeConfig(kc string) (string, *kubernetes.Clientset, error) {
	bs, err := url.PathUnescape(kc)
	if err != nil {
		return "", nil, err
	}
	ac, err := clientcmd.Load([]byte(bs))
	if err != nil {
		return "", nil, err
	}

	var ns = ""
	for _, ctx := range ac.Contexts {
		ns = ctx.Namespace
	}

	if ns == "" {
		return ns, nil, errors.New("认证信息错误，未发现命名空间")
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(bs))

	if err != nil {
		return ns, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ns, nil, err
	}
	return ns, client, nil
}
