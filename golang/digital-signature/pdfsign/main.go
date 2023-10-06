package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/digitorus/pdfsign/sign"
	"github.com/digitorus/pdfsign/verify"
)

var now = time.Now()

func generateKeys() (*rsa.PrivateKey, *x509.Certificate, error) {
	// Generate private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Initialize x509 certificate template
	template := x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   "any",
			Organization: []string{"Test Company"},
		},
		NotBefore:          now.Add(-time.Hour).UTC(),
		NotAfter:           now.Add(time.Hour * 24 * 365).UTC(),
		PublicKeyAlgorithm: x509.RSA,
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
	}

	// Generate x509 certificate
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

func usage() {
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Printf("\t%s sign input.pdf output.pdf\n", os.Args[0])
	fmt.Printf("\t%s verify input.pdf\n", os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
	}

	method := flag.Arg(0)
	if method != "sign" && method != "verify" {
		usage()
	}

	input := flag.Arg(1)
	if len(input) == 0 {
		usage()
	}

	if method == "sign" {
		output := flag.Arg(2)
		if len(output) == 0 {
			usage()
		}

		// Generate keys
		priv, cert, err := generateKeys()
		if err != nil {
			panic(err)
		}

		if err = sign.SignFile(input, output, sign.SignData{
			Signature: sign.SignDataSignature{
				Info: sign.SignDataSignatureInfo{
					Name:        "Kien",
					Location:    "Hanoi",
					Reason:      "Thich thi ky thoi",
					ContactInfo: "callmemaybe@mail.com",
					Date:        time.Now().Local(),
				},
				CertType:   sign.CertificationSignature,
				DocMDPPerm: sign.AllowFillingExistingFormFieldsAndSignaturesAndCRUDAnnotationsPerms,
			},
			Signer:          priv,
			DigestAlgorithm: crypto.SHA256,
			Certificate:     cert,
		}); err != nil {
			panic(err)
		}

		fmt.Println("Signed PDF written to ", output)
	}

	if method == "verify" {
		input_file, err := os.Open(input)
		if err != nil {
			panic(err)
		}
		defer input_file.Close()

		resp, err := verify.File(input_file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
		return
	}
}
