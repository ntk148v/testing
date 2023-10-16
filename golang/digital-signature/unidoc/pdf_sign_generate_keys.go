package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

var now = time.Now()

const usagef = "Usage: %s INPUT_PDF_PATH OUTPUT_PDF_PATH\n"

func init() {
	// // Make sure to load your metered License API key prior to using the library.
	// // If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	// err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	// if err != nil {
	// 	panic(err)
	// }
}

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

func exportRsaPrivateKeyasPemStr(priv *rsa.PrivateKey) string {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(priv)
	privKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA",
		Bytes: privKeyBytes,
	})

	return string(privKeyPem)
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf(usagef, os.Args[0])
		return
	}

	inputPath := args[1]
	outputPath := args[2]

	// Generate key pair
	priv, cert, err := generateKeys()
	if err != nil {
		panic(err)
	}

	// Create reader
	f, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		panic(err)
	}

	// Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		panic(err)
	}

	// Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		panic(err)
	}

	// Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName("Test Self Signed PDF")
	signature.SetReason("TestSelfSignedPDF")
	signature.SetDate(now, "")

	if err := signature.Initialize(); err != nil {
		panic(err)
	}

	// Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = 10
	opts.Rect = []float64{10, 25, 75, 60}

	field, _ := annotator.NewSignatureField(
		signature,
		[]*annotator.SignatureLine{
			annotator.NewSignatureLine("Name", "John Doe"),
			annotator.NewSignatureLine("Date", "2019.16.04"),
			annotator.NewSignatureLine("Reason", "External signature test"),
		},
		opts,
	)
	field.T = core.MakeString("Self signed PDF")

	if err = appender.Sign(1, field); err != nil {
		panic(err)
	}

	// Write output PDF file.
	err = appender.WriteToFile(outputPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
}
