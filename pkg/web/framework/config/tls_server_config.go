package config

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"

	"github.com/secretflow/kuscia/pkg/utils/network"
	"github.com/secretflow/kuscia/pkg/utils/nlog"
	"github.com/secretflow/kuscia/pkg/utils/tls"

	"github.com/google/uuid"
)

type TLSServerConfig struct {
	ServerCertFile string `yaml:"serverCertFile,omitempty"`
	ServerCertData string `yaml:"serverCertData,omitempty"`
	ServerKeyFile  string `yaml:"serverKeyFile,omitempty"`
	ServerKeyData  string `yaml:"serverKeyData,omitempty"`

	RootCA     *x509.Certificate `yaml:"-"`
	RootCAKey  *rsa.PrivateKey   `yaml:"-"`
	ServerCert *x509.Certificate `yaml:"-"`
	ServerKey  *rsa.PrivateKey   `yaml:"-"`
}

func (t *TLSServerConfig) LoadFromDataOrFile() error {
	var err error
	if t.ServerKeyData != "" || t.ServerKeyFile != "" {
		if t.ServerKey, err = tls.ParseKey([]byte(t.ServerKeyData),
			t.ServerKeyFile); err != nil {
			return err
		}
	}

	if t.ServerCertData != "" || t.ServerCertFile != "" {
		if t.ServerCert, err = tls.ParseCert([]byte(t.ServerCertData),
			t.ServerCertFile); err != nil {
			return err
		}
	}

	return nil
}

func (t *TLSServerConfig) GenerateServerKeyCerts(commonName string, ipList []string, dnsList []string) error {
	ips := make([]net.IP, 0)
	for _, ipStr := range ipList {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			nlog.Warnf("Generate server key certs, found ip[%s] parse failed, skip", ipStr)
			continue
		}
		ips = append(ips, ip)
	}
	ips = append(ips, net.ParseIP("127.0.0.1"))
	hostIP, err := network.GetHostIP()
	if err != nil {
		nlog.Warnf("GenerateServerKeyCerts inject host ip failed: %s, skip", err.Error())
	} else {
		ips = append(ips, net.ParseIP(hostIP))
	}

	certTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(int64(uuid.New().ID())),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		IPAddresses: ips,
		DNSNames:    dnsList,
	}
	t.ServerKey, t.ServerCert, err = tls.GenerateX509KeyPairStruct(t.RootCA, t.RootCAKey, certTmpl)
	if err != nil {
		return err
	}
	return nil
}
