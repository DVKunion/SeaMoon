package sealos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

var (
	num    int32 = 1
	prefix       = networkingv1.PathTypePrefix
	fl           = false
)

type amount struct {
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

func getAmountAndCost(ca *models.CloudAuth, region string) (float64, float64, error) {
	amountUrl := fmt.Sprintf("https://costcenter.%s/api/account/getAmount", regionMap[region])

	req, err := http.NewRequest("GET", amountUrl, nil)
	if err != nil {
		return 0, 0, err
	}

	req.Header.Add("Authorization", url.PathEscape(ca.KubeConfig))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, errors.New("error request : " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var sa = amount{}
	err = json.Unmarshal(body, &sa)
	if err != nil {
		return 0, 0, err
	}

	if sa.Code != 200 {
		return 0, 0, errors.New(sa.Message)
	}

	return float64(sa.Data.Balance-sa.Data.DeductionBalance) / 1000000, float64(sa.Data.DeductionBalance) / 1000000, nil
}

func deploy(config, svcName, imgName, hostName string, port int32, tc *models.TunnelConfig, tp *enum.TunnelType) (string, error) {
	ctx := context.Background()
	uid := ""
	ns, clientSet, err := parseKubeConfig(config)

	if err != nil {
		return "", err
	}

	if res, err := clientSet.AppsV1().Deployments(ns).
		Create(ctx, renderDeployment(svcName, imgName, port, tc, tp),
			metav1.CreateOptions{}); err != nil {
		return "", err
	} else {
		uid = string(res.ObjectMeta.UID)
	}

	if _, err = clientSet.CoreV1().Services(ns).
		Create(ctx, renderService(svcName, port), metav1.CreateOptions{}); err != nil {
		return "", err
	}

	// ingress
	if _, err = clientSet.NetworkingV1().Ingresses(ns).
		Create(ctx, renderIngress(svcName, hostName, tc, tp), metav1.CreateOptions{}); err != nil {
		return "", err
	}

	cnt := 0
	status := enum.TunnelInitializing
	message := ""

	for cnt < 30 {
		// 查看一下状态：
		svcs, err := clientSet.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{
			LabelSelector: "cloud.sealos.io/app-deploy-manager=" + svcName,
		})
		if err != nil {
			return "", err
		}

		if len(svcs.Items) == 0 {
			return "", errors.New(xlog.SDKFCCreateError)
		}

		for _, svc := range svcs.Items {
			if svc.ObjectMeta.Name == svcName {
				for _, condition := range svc.Status.Conditions {
					xlog.Info(xlog.SDKWaitingFCStatus, "type", condition.Type, "status", condition.Status, "cnt", cnt)
					if condition.Type == "Available" && condition.Status == "True" {
						status = enum.TunnelActive
						message = ""
						cnt = 31
						break
					}
					if condition.Type == "Progressing" && condition.Status == "True" {
						message = condition.Message
					}
					if condition.Type == "Progressing" && condition.Status == "False" {
						message = condition.Message
					}
					if condition.Type == "Available" && condition.Status == "False" && message == "" {
						status = enum.TunnelError
						message = condition.Message
					}
				}
			}
		}

		cnt += 1
		time.Sleep(2 * time.Second)
	}

	if status != enum.TunnelActive {
		return "", errors.New(message)
	}

	return uid, nil
}

func destroy(config, svcName string) error {
	ctx := context.Background()

	ns, clientSet, err := parseKubeConfig(config)

	if err != nil {
		return err
	}

	if err = clientSet.AppsV1().Deployments(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err = clientSet.CoreV1().Services(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err = clientSet.NetworkingV1().Ingresses(ns).Delete(ctx, svcName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func sync(config string) (*appsv1.DeploymentList, *networkingv1.IngressList, error) {
	ctx := context.Background()

	ns, clientSet, err := parseKubeConfig(config)

	if err != nil {
		return nil, nil, err
	}

	svcs, err := clientSet.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	ingresses, err := clientSet.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	return svcs, ingresses, nil
}

func renderDeployment(svcName, imgName string, port int32, config *models.TunnelConfig, tp *enum.TunnelType) *appsv1.Deployment {
	return &appsv1.Deployment{
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
								if config.Tor {
									// 需要增加环境变量
									env = append(env, corev1.EnvVar{
										Name:  "SEAMOON_TOR",
										Value: "true",
									})
								}
								env = append(env, corev1.EnvVar{
									Name:  "SM_UID",
									Value: config.V2rayUid,
								})
								env = append(env, corev1.EnvVar{
									Name:  "SM_SS_CRYPT",
									Value: config.SSRCrypt,
								})
								env = append(env, corev1.EnvVar{
									Name:  "SM_SS_PASS",
									Value: config.SSRPass,
								})
								return env
							}(),
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU: func() resource.Quantity {
										if config.CPU < 0.1 {
											return resource.MustParse("10m")
										}
										return resource.MustParse(strconv.Itoa(int(config.CPU*100)) + "m")
									}(),
									corev1.ResourceMemory: func() resource.Quantity {
										if config.Memory < 64 {
											return resource.MustParse("6Mi")
										}
										return resource.MustParse(strconv.Itoa(int(config.Memory/10)) + "Mi")
									}(),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU: func() resource.Quantity {
										if config.CPU < 0.1 {
											return resource.MustParse("100m")
										}
										return resource.MustParse(strconv.Itoa(int(config.CPU*1000)) + "m")
									}(),
									corev1.ResourceMemory: func() resource.Quantity {
										if config.Memory < 64 {
											return resource.MustParse("64Mi")
										}
										return resource.MustParse(strconv.Itoa(int(config.Memory)) + "Mi")
									}(),
								},
							},
							Args: func() []string {
								switch *tp {
								case enum.TunnelTypeWST:
									return []string{"server", "-p", "9000", "-t", "websocket"}
								case enum.TunnelTypeGRT:
									return []string{"server", "-p", "8089", "-t", "grpc"}
								}
								return []string{}
							}(),
							Ports: []corev1.ContainerPort{
								{
									Name:          "seamoon-http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: port,
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}

}

func renderService(svcName string, port int32) *corev1.Service {
	return &corev1.Service{
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
}

func renderIngress(svcName, hostName string, config *models.TunnelConfig, tp *enum.TunnelType) *networkingv1.Ingress {
	return &networkingv1.Ingress{
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
					switch *tp {
					case enum.TunnelTypeWST:
						return "WS"
					case enum.TunnelTypeGRT:
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
					Host: fmt.Sprintf("%s.%s", hostName, regionMap[config.Region]),
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
					Hosts:      []string{fmt.Sprintf("%s.%s", hostName, regionMap[config.Region])},
					SecretName: "wildcard-cloud-sealos-io-cert",
				},
			},
		},
	}
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
