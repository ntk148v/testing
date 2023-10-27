package org.example;

import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.model.*;
import eu.europa.esig.dss.pades.PAdESSignatureParameters;
import eu.europa.esig.dss.pades.SignatureFieldParameters;
import eu.europa.esig.dss.pades.SignatureImageParameters;
import eu.europa.esig.dss.pades.signature.PAdESService;
import eu.europa.esig.dss.pdf.pdfbox.PdfBoxNativeObjectFactory;
import eu.europa.esig.dss.pdfa.PDFAValidationResult;
import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
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

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.Certificate;
import java.security.spec.RSAKeyGenParameterSpec;
import java.util.Collection;
import java.util.Date;
import java.util.List;

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