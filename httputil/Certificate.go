package httputil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func CertFileCheck(cert_file string) (output bool, err error) {
	if _, err := os.Stat(cert_file); os.IsNotExist(err) {
		output = false
	} else {
		output = true
	}
	return
}
func Certify(cert_file string) {
	check, _ := CertFileCheck(cert_file)
	if !check {
		max := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, _ := rand.Int(rand.Reader, max)

		subject := pkix.Name{
			Organization:       []string{"Certificate File"},
			OrganizationalUnit: []string{"File"},
			CommonName:         "Glue API",
		}
		template := x509.Certificate{
			SerialNumber: serialNumber,
			Subject:      subject,
			NotBefore:    time.Now(),
			NotAfter:     time.Now().Add(365 * 24 * time.Hour),
			KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		}
		pk, _ := rsa.GenerateKey(rand.Reader, 2048)

		derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
		certOut, _ := os.Create("cert.pem")
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		certOut.Close()

		keyOut, _ := os.Create("key.pem")
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
		keyOut.Close()
	}

}
