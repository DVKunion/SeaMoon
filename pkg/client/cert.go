package client

import (
	"bytes"
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"os"
	"time"
)

var CA = "ca.crt"
var Key = "ca.key"

type CERT struct {
	CERT       []byte
	CERTKEY    *rsa.PrivateKey
	CERTPEM    *bytes.Buffer
	CERTKEYPEM *bytes.Buffer
	CSR        *x509.Certificate
}

func InitCa() error {
	_, errCa := os.Stat(CA)
	_, errKey := os.Stat(Key)
	if os.IsNotExist(errCa) || os.IsNotExist(errKey) {
		log.Info(consts.CA_NOT_EXIST)
		err := generateCa()
		if err != nil {
			return err
		}
	}
	caCert, err := ioutil.ReadFile(CA)
	if err != nil {
		return err
	}

	caKey, err := ioutil.ReadFile(Key)
	if err != nil {
		return err
	}

	err = setCA(caCert, caKey)
	if err != nil {
		return err
	}
	log.Info(consts.CA_LOAD_SUCCESS)
	return nil
}

func setCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

func generateCa() error {
	subj := &pkix.Name{
		CommonName:    "SeaMoon",
		Organization:  []string{"Company, INC."},
		Country:       []string{"CH"},
		Province:      []string{""},
		Locality:      []string{"HangZhou"},
		StreetAddress: []string{"HangZhou"},
		PostalCode:    []string{"94016"},
	}
	ca, err := create(subj, 10)
	if err != nil {
		return err
	}
	write(ca, "./ca")
	crt, err := req(ca.CSR, subj, 10)
	if err != nil {
		return err
	}
	write(crt, "./tls")
	return nil
}

func create(sub *pkix.Name, expire int) (*CERT, error) {
	var (
		ca  = new(CERT)
		err error
	)

	if expire < 1 {
		expire = 1
	}
	// 为ca生成私钥
	ca.CERTKEY, err = rsa.GenerateKey(cr.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// 对证书进行签名
	ca.CSR = &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63n(2000)),
		Subject:      *sub,
		NotBefore:    time.Now(),                       // 生效时间
		NotAfter:     time.Now().AddDate(expire, 0, 0), // 过期时间
		IsCA:         true,                             // 表示用于CA
		// openssl 中的 extendedKeyUsage = clientAuth, serverAuth 字段
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		// openssl 中的 keyUsage 字段
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// 创建证书
	// caBytes 就是生成的证书
	ca.CERT, err = x509.CreateCertificate(cr.Reader, ca.CSR, ca.CSR, &ca.CERTKEY.PublicKey, ca.CERTKEY)
	if err != nil {
		return nil, err
	}
	ca.CERTPEM = new(bytes.Buffer)
	pem.Encode(ca.CERTPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.CERT,
	})
	ca.CERTKEYPEM = new(bytes.Buffer)
	pem.Encode(ca.CERTKEYPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca.CERTKEY),
	})

	// 进行PEM编码，编码就是直接cat证书里面内容显示的东西
	return ca, nil
}

func req(ca *x509.Certificate, sub *pkix.Name, expire int) (*CERT, error) {
	var (
		cert = &CERT{}
		err  error
	)
	cert.CERTKEY, err = rsa.GenerateKey(cr.Reader, 4096)
	if err != nil {
		return nil, err
	}
	if expire < 1 {
		expire = 1
	}
	cert.CSR = &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63n(2000)),
		Subject:      *sub,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(expire, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	cert.CERT, err = x509.CreateCertificate(cr.Reader, cert.CSR, ca, &cert.CERTKEY.PublicKey, cert.CERTKEY)
	if err != nil {
		return nil, err
	}

	cert.CERTPEM = new(bytes.Buffer)
	pem.Encode(cert.CERTPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.CERT,
	})
	cert.CERTKEYPEM = new(bytes.Buffer)
	pem.Encode(cert.CERTKEYPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.CERTKEY),
	})
	return cert, nil
}

func write(cert *CERT, file string) error {
	keyFileName := file + ".key"
	certFIleName := file + ".crt"
	kf, err := os.Create(keyFileName)
	if err != nil {
		return err
	}
	defer kf.Close()

	if _, err := kf.Write(cert.CERTKEYPEM.Bytes()); err != nil {
		return err
	}

	cf, err := os.Create(certFIleName)
	if err != nil {
		return err
	}
	if _, err := cf.Write(cert.CERTPEM.Bytes()); err != nil {
		return err
	}
	return nil
}

func load(caFile string) (cp *x509.CertPool, err error) {
	if caFile == "" {
		return
	}
	cp = x509.NewCertPool()
	data, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	if !cp.AppendCertsFromPEM(data) {
		return nil, errors.New("AppendCertsFromPEM failed")
	}
	return
}

func GetClientTLSConfig(addr, ca string, skipVerify bool) (config *tls.Config, err error) {
	rootCAs, err := load(ca)
	if err != nil {
		return
	}
	serverName, _, _ := net.SplitHostPort(addr)
	if net.ParseIP(serverName) != nil { // server name is IP
		config = &tls.Config{
			InsecureSkipVerify: true,
			VerifyConnection: func(cs tls.ConnectionState) error { // verify manually
				if skipVerify {
					return nil
				}

				opts := x509.VerifyOptions{
					Roots:         rootCAs,
					CurrentTime:   time.Now(),
					Intermediates: x509.NewCertPool(),
				}

				certs := cs.PeerCertificates
				for i, cert := range certs {
					if i == 0 {
						continue
					}
					opts.Intermediates.AddCert(cert)
				}

				_, err := certs[0].Verify(opts)
				return err
			},
		}
	} else { // server name is domain
		config = &tls.Config{
			ServerName:         serverName,
			RootCAs:            rootCAs,
			InsecureSkipVerify: skipVerify,
		}
	}

	return
}
