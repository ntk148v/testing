package org.example;

import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.model.*;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.pades.PAdESSignatureParameters;
import eu.europa.esig.dss.pades.SignatureFieldParameters;
import eu.europa.esig.dss.pades.SignatureImageParameters;
import eu.europa.esig.dss.pades.signature.PAdESService;
import eu.europa.esig.dss.pdf.pdfbox.PdfBoxNativeObjectFactory;
import eu.europa.esig.dss.pdfa.PDFAValidationResult;
import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
import eu.europa.esig.dss.token.MSCAPISignatureToken;
import eu.europa.esig.dss.token.Pkcs11SignatureToken;
import eu.europa.esig.dss.pdfa.validation.PDFADocumentValidator;
import eu.europa.esig.dss.utils.Utils;
import eu.europa.esig.dss.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.reports.Reports;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

import java.io.*;
import java.lang.reflect.Array;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.spec.RSAKeyGenParameterSpec;
import java.util.*;

public class Main {
    private static String HSM_PIN = "1234";
    private static final int KEY_SIZE = 2048;
    private static String KEY_PASSWD = "1234";
    private static String KEY_ALIAS = "keycuakien";

    private static void genAndAddKeyStore() {
        try {
            String configName = "softhsm.cfg";
            Provider sunPkcsProvider = Security.getProvider("SunPKCS11");
            sunPkcsProvider = sunPkcsProvider.configure(configName);
            Security.addProvider(sunPkcsProvider);

            KeyStore keyStore = KeyStore.getInstance("PKCS11", sunPkcsProvider);
            keyStore.load(null, HSM_PIN.toCharArray());

            // Create keypair
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            SecureRandom random = new SecureRandom();
            keyPairGenerator.initialize(new RSAKeyGenParameterSpec(KEY_SIZE, RSAKeyGenParameterSpec.F4), random);
            KeyPair keyPair = keyPairGenerator.generateKeyPair();
            keyPair.getPrivate();

            X509v3CertificateBuilder certBuilder = new JcaX509v3CertificateBuilder(new X500Name("CN=Test"), BigInteger.valueOf(1), // serial
                    new Date(System.currentTimeMillis() - 1000L * 60 * 60 * 24 * 30), new Date(System.currentTimeMillis() + (1000L * 60 * 60 * 24 * 30)), new X500Name("CN=Test"), keyPair.getPublic());
            ContentSigner contentSigner;
            contentSigner = new JcaContentSignerBuilder("SHA256WithRSAEncryption").build(keyPair.getPrivate());
            Certificate certificate = new JcaX509CertificateConverter().setProvider(new BouncyCastleProvider()).getCertificate(certBuilder.build(contentSigner));

            keyStore.setKeyEntry(KEY_ALIAS, // alias
                    keyPair.getPrivate(), // privatekey
                    KEY_PASSWD.toCharArray(), // keystore password to protect key
                    new Certificate[]{certificate}); // the cert chain for its public key
            keyStore.store(null, HSM_PIN.toCharArray());

            System.out.println("Added key");
            System.out.println("Get private to use");
            PrivateKey privateKey = (PrivateKey) keyStore.getKey(KEY_ALIAS, KEY_PASSWD.toCharArray());
            // Attempt to print private key
            //            The HSM will let you perform crypto operations using the stored private key,
            //            but won’t let you access the content of the key itself. Here is what we get if we
            //            try to print the key to the console via the Java code.
            // SunPKCS11-SoftHSM RSA private key, 2048 bits token object, not sensitive, unextractable)
            System.out.println(privateKey);


        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public static void main(String[] args) {
        String certStr = "MIIFPTCCBCWgAwIBAgIQVAEBYZ82L1CjhIy7Ds8Z6TANBgkqhkiG9w0BAQsFADBbMRcwFQYDVQQDDA5GYXN0Q0EgU0hBLTI1NjEzMDEGA1UECgwqQ8OUTkcgVFkgQ+G7lCBQSOG6pk4gQ0jhu64gS8OdIFPhu5AgRkFTVENBMQswCQYDVQQGEwJWTjAeFw0yMzExMDMwODEyMDBaFw0yNDExMDMwODEyMDBaMGQxCzAJBgNVBAYTAlZOMRIwEAYDVQQIDAlIw6AgTuG7mWkxHjAcBgNVBAMMFU5HVVnhu4ROIFRV4bqkTiBLScOKTjEhMB8GCgmSJomT8ixkAQEMEUNDQ0Q6Mjg2MDk0MDAwMDIwMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2KcyXepH6YS25XoKR6yZxTxpheKJhtkljqKyW3yIYXFaC6Iwe1eQmOEcdng1TYYl7o8lLaNe2gWn10d2SScI9xkvnNCGmTVuZ0W0LUPwnz2TiUCcsdr9YNdnltlXofPV0+0NG5kmmf3e4+AYudmTUIyXyx1cu9c+YZrwS7j5CXktYaSoM81GL9+gUV27SGHuZXOYvOynxJeYJP18wMxouZESTZtxqPDS8PjtPuj4BTs8DUvtpnllBv/avT6CEYAxPn4aEATi3iEO+eehEI8H5aWUSYnsG97bgeeTB64s0g/Dwat183wn/zmrVFWbEbe8gt0QHI4UtvrSRworf7QeewIDAQABo4IB8jCCAe4wDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBRSKOOvdOPsfmMJUHch3ZqOGP1KfDBeBggrBgEFBQcBAQRSMFAwKwYIKwYBBQUHMAKGH2h0dHA6Ly9wdWIuZmFzdGNhLnZuL0Zhc3RDQS5jcnQwIQYIKwYBBQUHMAGGFWh0dHA6Ly9vY3NwLmZhc3RjYS52bjBhBgNVHSAEWjBYMFYGDCsGAQQBge0DAQkDATBGMBYGCCsGAQUFBwICMAoMCFBlcnNvbmFsMCwGCCsGAQUFBwICMCAMHmh0dHBzOi8vZmFzdGNhLnZuL2Rvd25sb2FkL0NQUzA0BgNVHSUELTArBggrBgEFBQcDAgYIKwYBBQUHAwQGCisGAQQBgjcKAwwGCSqGSIb3LwEBBTCBlAYDVR0fBIGMMIGJMIGGoCOgIYYfaHR0cDovL2NybC5mYXN0Y2Eudm4vZmFzdGNhLmNybKJfpF0wWzEXMBUGA1UEAwwORmFzdENBIFNIQS0yNTYxMzAxBgNVBAoMKkPDlE5HIFRZIEPhu5QgUEjhuqZOIENI4buuIEvDnSBT4buQIEZBU1RDQTELMAkGA1UEBhMCVk4wHQYDVR0OBBYEFNg29r6rfh8pzd0V7NGH93qPH2VjMA4GA1UdDwEB/wQEAwIE8DANBgkqhkiG9w0BAQsFAAOCAQEAtPzyvc1QAnW5EztuImg0EcKmklhvear9rG20SmnkTqR7fsQkLm5ext5uMtPF1OucrQ8AoVT8Dd1TwOzxpiv5JRUaanFo1b7aet7Ql2wjXsyKt29uGinWGZSQFV/xa7QdFp5uqjg3yhRoN5CfmtrB6lJDuHc1MibKmgChEQHlaoQBczS2AfXUQgN30NwUIy0ujSDOXIXiWbyS9Ijt/Eck1a2jBNlnA91mdJsCjWevgUBiAQ1fNtnkxyt7a64LDLcvAlpOFgvyDoc/F4wGnGLOgRZT5OUGay+fbJAjiPJl/T6PLt2FhxyhgD4+FRRsl5LQS/NMYUQ0g3UPDy9zPSVEkA==";
        List<String> certChainStr = Arrays.asList(
                "MIIFPTCCBCWgAwIBAgIQVAEBYZ82L1CjhIy7Ds8Z6TANBgkqhkiG9w0BAQsFADBbMRcwFQYDVQQDDA5GYXN0Q0EgU0hBLTI1NjEzMDEGA1UECgwqQ8OUTkcgVFkgQ+G7lCBQSOG6pk4gQ0jhu64gS8OdIFPhu5AgRkFTVENBMQswCQYDVQQGEwJWTjAeFw0yMzExMDMwODEyMDBaFw0yNDExMDMwODEyMDBaMGQxCzAJBgNVBAYTAlZOMRIwEAYDVQQIDAlIw6AgTuG7mWkxHjAcBgNVBAMMFU5HVVnhu4ROIFRV4bqkTiBLScOKTjEhMB8GCgmSJomT8ixkAQEMEUNDQ0Q6Mjg2MDk0MDAwMDIwMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2KcyXepH6YS25XoKR6yZxTxpheKJhtkljqKyW3yIYXFaC6Iwe1eQmOEcdng1TYYl7o8lLaNe2gWn10d2SScI9xkvnNCGmTVuZ0W0LUPwnz2TiUCcsdr9YNdnltlXofPV0+0NG5kmmf3e4+AYudmTUIyXyx1cu9c+YZrwS7j5CXktYaSoM81GL9+gUV27SGHuZXOYvOynxJeYJP18wMxouZESTZtxqPDS8PjtPuj4BTs8DUvtpnllBv/avT6CEYAxPn4aEATi3iEO+eehEI8H5aWUSYnsG97bgeeTB64s0g/Dwat183wn/zmrVFWbEbe8gt0QHI4UtvrSRworf7QeewIDAQABo4IB8jCCAe4wDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBRSKOOvdOPsfmMJUHch3ZqOGP1KfDBeBggrBgEFBQcBAQRSMFAwKwYIKwYBBQUHMAKGH2h0dHA6Ly9wdWIuZmFzdGNhLnZuL0Zhc3RDQS5jcnQwIQYIKwYBBQUHMAGGFWh0dHA6Ly9vY3NwLmZhc3RjYS52bjBhBgNVHSAEWjBYMFYGDCsGAQQBge0DAQkDATBGMBYGCCsGAQUFBwICMAoMCFBlcnNvbmFsMCwGCCsGAQUFBwICMCAMHmh0dHBzOi8vZmFzdGNhLnZuL2Rvd25sb2FkL0NQUzA0BgNVHSUELTArBggrBgEFBQcDAgYIKwYBBQUHAwQGCisGAQQBgjcKAwwGCSqGSIb3LwEBBTCBlAYDVR0fBIGMMIGJMIGGoCOgIYYfaHR0cDovL2NybC5mYXN0Y2Eudm4vZmFzdGNhLmNybKJfpF0wWzEXMBUGA1UEAwwORmFzdENBIFNIQS0yNTYxMzAxBgNVBAoMKkPDlE5HIFRZIEPhu5QgUEjhuqZOIENI4buuIEvDnSBT4buQIEZBU1RDQTELMAkGA1UEBhMCVk4wHQYDVR0OBBYEFNg29r6rfh8pzd0V7NGH93qPH2VjMA4GA1UdDwEB/wQEAwIE8DANBgkqhkiG9w0BAQsFAAOCAQEAtPzyvc1QAnW5EztuImg0EcKmklhvear9rG20SmnkTqR7fsQkLm5ext5uMtPF1OucrQ8AoVT8Dd1TwOzxpiv5JRUaanFo1b7aet7Ql2wjXsyKt29uGinWGZSQFV/xa7QdFp5uqjg3yhRoN5CfmtrB6lJDuHc1MibKmgChEQHlaoQBczS2AfXUQgN30NwUIy0ujSDOXIXiWbyS9Ijt/Eck1a2jBNlnA91mdJsCjWevgUBiAQ1fNtnkxyt7a64LDLcvAlpOFgvyDoc/F4wGnGLOgRZT5OUGay+fbJAjiPJl/T6PLt2FhxyhgD4+FRRsl5LQS/NMYUQ0g3UPDy9zPSVEkA==",
                "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",
                "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"
        );
        System.out.println(LoadCertificate(certStr));
        for (String cs : certChainStr) {
            System.out.println(LoadCertificate(cs));
        }
    }

    public static CertificateToken LoadCertificate(String certStr) {
        X509Certificate cert = null;
        Base64.Decoder decoder = Base64.getDecoder();
        byte[] decodedData = decoder.decode(certStr);
        try (InputStream inputStream = new ByteArrayInputStream(decodedData)) {
            CertificateFactory cf = CertificateFactory.getInstance("X.509");

            java.security.cert.Certificate certificate = cf.generateCertificate(inputStream);

            if (certificate instanceof X509Certificate) {
                cert = (X509Certificate) certificate;
                return new CertificateToken(cert);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

        return null;
    }

    public static void MSCAPISign() {
        try (MSCAPISignatureToken token = new MSCAPISignatureToken()) {
            DSSPrivateKeyEntry privateKey = token.getKeys().get(0);
            System.out.println(privateKey.getCertificate());
//            ToBeSigned toBeSigned = new eu.europa.esig.dss.model.ToBeSigned("Hello world".getBytes());
//            SignatureValue signatureValue = token.sign(toBeSigned, DigestAlgorithm.SHA256, privateKey);
//
//            System.out.println("Signature value : " + Utils.toBase64(signatureValue.getValue()));
            // Sign a PDF
            // Preparing parameters for the PAdES signature
            PAdESSignatureParameters parameters = new PAdESSignatureParameters();
            // We choose the level of the signature (-B, -T, -LT, -LTA).
            parameters.setSignatureLevel(SignatureLevel.PAdES_BASELINE_B);
            // We set the digest algorithm to use with the signature algorithm. You must use the
            // same parameter when you invoke the method sign on the token. The default value is
            // SHA256
            parameters.setDigestAlgorithm(DigestAlgorithm.SHA256);

            // We set the signing certificate
            parameters.setSigningCertificate(privateKey.getCertificate());
            // Create common certificate verifier
            CommonCertificateVerifier commonCertificateVerifier = new CommonCertificateVerifier();
            // Create PAdESService for signature
            PAdESService service = new PAdESService(commonCertificateVerifier);
            service.setPdfObjFactory(new PdfBoxNativeObjectFactory());

            // Do it in the web client, return the data!
            DSSDocument toSignDocument = new FileDocument("input.pdf");
            // Get the SignedInfo segment that need to be signed.
            ToBeSigned dataToSign = service.getDataToSign(toSignDocument, parameters);
            System.out.println("To be signed value : " + Utils.toBase64(dataToSign.getBytes()));
            ;
            // This function obtains the signature value for signed information using the
            // private key and specified algorithm
            DigestAlgorithm digestAlgorithm = DigestAlgorithm.SHA256;
            // Do it in the client application
//            SignatureValue signatureValue = token.sign(dataToSign, digestAlgorithm, privateKey);
//            System.out.println("Signature value : " + Utils.toBase64(signatureValue.getValue()));
//
//            DSSDocument signedDocument = service.signDocument(toSignDocument, parameters, signatureValue);
//            signedDocument.writeTo(new FileOutputStream("output2.pdf"));
        }
    }

    public static void PKCS11Sign() {
        // uncomment if you want to gen and add :x
        // genAndAddKeyStore();
        try (Pkcs11SignatureToken token = new Pkcs11SignatureToken("/usr/local/lib/softhsm/libsofthsm2.so", new KeyStore.PasswordProtection("1234".toCharArray()), 1651866022)) {
            List<DSSPrivateKeyEntry> keys = token.getKeys();
//            for (DSSPrivateKeyEntry entry : keys) {
//                System.out.println(entry.getCertificate().getCertificate());
//                System.out.println(entry.getCertificate().isSelfSigned());
//            }

            DSSPrivateKeyEntry privateKey = token.getKeys().get(0);
            System.out.printf("I got the privatekey entry it's self-signed = %s\n", privateKey.getCertificate().isSelfSigned());

            // Sign a byte
//            ToBeSigned toBeSigned = new ToBeSigned("Hello world".getBytes());
//            System.out.println("Sign ne");
//            SignatureValue signatureValue = token.sign(toBeSigned, DigestAlgorithm.SHA256, privateKey); // use the first key only
//            System.out.println("Signature value : " + Utils.toBase64(signatureValue.getValue()));

            // Sign a PDF
            // Preparing parameters for the PAdES signature
            PAdESSignatureParameters parameters = new PAdESSignatureParameters();
            // We choose the level of the signature (-B, -T, -LT, -LTA).
            parameters.setSignatureLevel(SignatureLevel.PAdES_BASELINE_B);
            // We set the digest algorithm to use with the signature algorithm. You must use the
            // same parameter when you invoke the method sign on the token. The default value is
            // SHA256
            parameters.setDigestAlgorithm(DigestAlgorithm.SHA256);

            // We set the signing certificate
            parameters.setSigningCertificate(privateKey.getCertificate());
            // We set the certificate chain
            parameters.setCertificateChain(privateKey.getCertificateChain());

            // Initialize visual signature and configure
            SignatureImageParameters imageParameters = new SignatureImageParameters();
            // set an image
            imageParameters.setImage(new InMemoryDocument(new FileInputStream("signed.jpeg")));

            // initialize signature field parameters
            SignatureFieldParameters fieldParameters = new SignatureFieldParameters();
            imageParameters.setFieldParameters(fieldParameters);
            // the origin is the left and top corner of the page
            fieldParameters.setOriginX(200);
            fieldParameters.setOriginY(400);
            fieldParameters.setWidth(300);
            fieldParameters.setHeight(200);
            parameters.setImageParameters(imageParameters);

            // Create common certificate verifier
            CommonCertificateVerifier commonCertificateVerifier = new CommonCertificateVerifier();
            // Create PAdESService for signature
            PAdESService service = new PAdESService(commonCertificateVerifier);
            service.setPdfObjFactory(new PdfBoxNativeObjectFactory());

            // Do it in the web client, return the data!
            DSSDocument toSignDocument = new FileDocument("input.pdf");
            // Get the SignedInfo segment that need to be signed.
            ToBeSigned dataToSign = service.getDataToSign(toSignDocument, parameters);
            // This function obtains the signature value for signed information using the
            // private key and specified algorithm
            DigestAlgorithm digestAlgorithm = parameters.getDigestAlgorithm();
            // Do it in the client application
            SignatureValue signatureValue = token.sign(dataToSign, digestAlgorithm, privateKey);

            // Optionally or for debug purpose :
            // Validate the signature value against the original dataToSign
            System.out.printf("Validate the signature value against the original dataToSign: %s\n",
                    service.isValidSignatureValue(dataToSign, signatureValue, privateKey.getCertificate()));

            // Do it in the server side
            System.out.println("Add the signature value to the real document, it's output.pdf, check it out");
            // We invoke the padesService to sign the document with the signature value obtained in
            // the previous step.
            DSSDocument signedDocument = service.signDocument(toSignDocument, parameters, signatureValue);
            signedDocument.writeTo(new FileOutputStream("output.pdf"));

            // Validator
            // Create a PDFADocumentValidator to perform validation against PDF/A specification
            PDFADocumentValidator documentValidator = new PDFADocumentValidator(signedDocument);
            // Extract PDF/A validation result
            // This report contains only validation of a document against PDF/A specification
            // and no signature validation process result
            PDFAValidationResult pdfaValidationResult = documentValidator.getPdfValidationResult();
            // This variable contains the name of the identified PDF/A profile (or closest if validation failed)
            String profileId = pdfaValidationResult.getProfileId();

            // Checks whether the PDF document is compliant to the identified PDF profile
            boolean compliant = pdfaValidationResult.isCompliant();

            // Returns the error messages occurred during the PDF/A verification
            Collection<String> errorMessages = pdfaValidationResult.getErrorMessages();

            // It is also possible to perform the signature validation process and extract the PDF/A validation result from DiagnosticData
            // Configure PDF/A document validator and perform validation of the document
            documentValidator.setCertificateVerifier(commonCertificateVerifier);
            Reports reports = documentValidator.validateDocument();

            // Extract the interested information from DiagnosticData
            DiagnosticData diagnosticData = reports.getDiagnosticData();
            profileId = diagnosticData.getPDFAProfileId();
            compliant = diagnosticData.isPDFACompliant();
            errorMessages = diagnosticData.getPDFAValidationErrors();
        } catch (FileNotFoundException e) {
            throw new RuntimeException(e);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}