package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"github.com/miekg/pkcs11"
)

var now = time.Now()

func generateKeys() (*rsa.PrivateKey, *x509.Certificate, error) {
	// Generate private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Initialize X509 certificate template.
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Company"},
		},
		NotBefore: now.Add(-time.Hour),
		NotAfter:  now.Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Generate X509 certificate.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	return priv, cert, nil
}

func main() {
	// Install softhsm2
	// run it: softhsm2-util --init-token --slot 0 --label test --pin 123456 --so-pin 123456
	p := pkcs11.New("/usr/local/lib/softhsm/libsofthsm2.so")
	err := p.Initialize()
	if err != nil {
		panic(err)
	}

	defer func() {
		p.Destroy()
		p.Finalize()
	}()

	slots, err := p.GetSlotList(true)
	if err != nil {
		panic(err)
	}

	fmt.Println(slots) // You should get 2 slots, the first one is the one we created before.

	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		panic(err)
	}
	defer p.CloseSession(session)

	err = p.Login(session, pkcs11.CKU_USER, "123456") // yeah, I fixed the pin here
	if err != nil {
		panic(err)
	}

	defer p.Logout(session)
	p.DigestInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256, nil)})
	hash, err := p.Digest(session, []byte("test"))
	if err != nil {
		panic(err)
	}

	for _, d := range hash {
		fmt.Printf("%x", d)
	}
}

func generateRSAKeyPair(p *pkcs11.Ctx, session pkcs11.SessionHandle, tokenLabel string, tokenPersistent bool) (pkcs11.ObjectHandle, pkcs11.ObjectHandle) {
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 1}),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
	}
	pbk, pvk, err := p.GenerateKeyPair(session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
		publicKeyTemplate, privateKeyTemplate)
	if err != nil {
		panic(err)
	}

	return pbk, pvk
}
