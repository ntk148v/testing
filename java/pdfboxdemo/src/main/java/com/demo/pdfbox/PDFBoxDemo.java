package com.demo.pdfbox;

import org.apache.pdfbox.Loader;
import org.apache.pdfbox.cos.COSDictionary;
import org.apache.pdfbox.examples.signature.CreateVisibleSignature2;
import org.apache.pdfbox.examples.signature.SigUtils;
import org.apache.pdfbox.examples.signature.cert.CertificateVerificationException;
import org.apache.pdfbox.examples.signature.cert.CertificateVerifier;
import org.apache.pdfbox.pdfparser.PDFParser;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.io.RandomAccessReadBufferedFile;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.PDPageContentStream;
import org.apache.pdfbox.pdmodel.encryption.SecurityProvider;
import org.apache.pdfbox.pdmodel.graphics.image.PDImageXObject;
import org.apache.pdfbox.pdmodel.interactive.digitalsignature.COSFilterInputStream;
import org.apache.pdfbox.pdmodel.interactive.digitalsignature.PDSignature;
import org.bouncycastle.asn1.cms.Attribute;
import org.bouncycastle.asn1.cms.CMSAttributes;
import org.bouncycastle.asn1.cms.Time;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.cms.*;
import org.bouncycastle.cms.jcajce.JcaSimpleSignerInfoVerifierBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.apache.pdfbox.examples.signature.CreateVisibleSignature;
import org.bouncycastle.tsp.TSPException;
import org.bouncycastle.tsp.TimeStampToken;
import org.bouncycastle.util.CollectionStore;
import org.bouncycastle.util.Store;

import java.awt.geom.Rectangle2D;
import java.io.*;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.*;
import java.security.cert.Certificate;
import java.security.spec.RSAKeyGenParameterSpec;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.Collection;
import java.util.Date;
import java.util.HashSet;

public class PDFBoxDemo {
    private static final int KEY_SIZE = 2048;
    private final KeyPair keyPair;
    private KeyStore keystore;
    private X509Certificate certificate;
    private static final String IN_DIR = "src/main/resources/com.demo.pdfbox/";
    private static final String OUT_DIR = "target/";
    private static final String STAMP_PATH = IN_DIR + "stamp.jpg";
    private final SimpleDateFormat sdf = new SimpleDateFormat("dd.MM.yyyy HH:mm:ss");


    public PDFBoxDemo() {
        this(KEY_SIZE);
    }

    public PDFBoxDemo(int keySize) {
        KeyPairGenerator keyPairGenerator;
        try {
            keyPairGenerator = KeyPairGenerator.getInstance("RSA");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("RSA algo not available", e);
        }
        SecureRandom random = new SecureRandom();
        try {
            keyPairGenerator.initialize(new RSAKeyGenParameterSpec(keySize, RSAKeyGenParameterSpec.F4), random);
        } catch (InvalidAlgorithmParameterException e) {
            throw new RuntimeException("unsupported key size: " + keySize);
        }
        this.keyPair = keyPairGenerator.generateKeyPair();
        this.generateCertificate();
    }

    /**
     * Gives back the RSA private key.
     *
     * @return
     */
    public PrivateKey getPrivate() {
        return this.keyPair.getPrivate();
    }

    public X509Certificate getCertificate() {
        if (null == this.certificate) {
            generateCertificate();
        }
        return this.certificate;
    }

    public KeyStore getKeystore() {
        if (this.keystore == null) {
            createKeyStore();
        }

        return this.keystore;
    }

    private void createKeyStore() {
        try {
            this.keystore = KeyStore.getInstance("PKCS12");
            this.keystore.load(null, null);
            // Load the private key and certificate
            this.keystore.setKeyEntry("alias", this.keyPair.getPrivate(), null, new Certificate[]{this.getCertificate()});
        } catch (KeyStoreException | CertificateException | IOException | NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
    }

    private void generateCertificate() {
        // Add Bouncy Castle as a Security Provider
        Security.addProvider(new BouncyCastleProvider());
        X509v3CertificateBuilder certBuilder = new JcaX509v3CertificateBuilder(new X500Name("CN=Test"), BigInteger.valueOf(1), // serial
                new Date(System.currentTimeMillis() - 1000L * 60 * 60 * 24 * 30), new Date(System.currentTimeMillis() + (1000L * 60 * 60 * 24 * 30)), new X500Name("CN=Test"), this.keyPair.getPublic());

        ContentSigner contentSigner;
        try {
            contentSigner = new JcaContentSignerBuilder("SHA256WithRSAEncryption").build(this.keyPair.getPrivate());
            this.certificate = new JcaX509CertificateConverter().setProvider("BC").getCertificate(certBuilder.build(contentSigner));
        } catch (OperatorCreationException e) {
            throw new RuntimeException(e);
        } catch (CertificateException e) {
            throw new RuntimeException(e);
        }
    }

    public void sign1() throws IOException, UnrecoverableKeyException, CertificateException, KeyStoreException, NoSuchAlgorithmException {
        // input pdf
        String inPath = IN_DIR + "input.pdf";
        File outFile;
        try (FileInputStream fis = new FileInputStream(STAMP_PATH)) {
            CreateVisibleSignature signing = new CreateVisibleSignature(this.getKeystore(), null);
            signing.setVisibleSignDesigner(inPath, 0, 0, -50, fis, 1);
            signing.setVisibleSignatureProperties("name", "location", "Security", 0, 1, true);
            outFile = new File(OUT_DIR + "output1.pdf");
            signing.signPDF(new File(inPath), outFile, null);
        }
    }

    public void sign2() {
        // input pdf
        String inPath = IN_DIR + "input.pdf";
        File documentFile = new File(inPath);
        File signedDocumentFile;
        signedDocumentFile = new File(OUT_DIR + "output2.pdf");
        try {
            CreateVisibleSignature2 signing = new CreateVisibleSignature2(this.getKeystore(), null);
            signing.setImageFile(new File(STAMP_PATH));
            Rectangle2D humanRect = new Rectangle2D.Float(100, 200, 150, 50);
            signing.setExternalSigning(false);
            signing.signPDF(documentFile, signedDocumentFile, humanRect, null);
        } catch (IOException | UnrecoverableKeyException | CertificateException | KeyStoreException |
                 NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
    }

    public void insertImage(String filename) throws IOException {
        PDDocument document = Loader.loadPDF(new File(OUT_DIR + filename));
        PDPage page = document.getPage(0);
        PDImageXObject pdImage = PDImageXObject.createFromFile(STAMP_PATH, document);
        PDPageContentStream contentStream = new PDPageContentStream(document, page);
        contentStream.drawImage(pdImage, 70, 250);
        contentStream.close();
        document.save(OUT_DIR + filename); // overwrite
        document.close();
    }

    public void verify(String filename) throws IOException {
        // output pdf
        // sign -> outFile -> become inFile of verify
        File outFile = new File(OUT_DIR + filename);
        // use old-style document loading to disable leniency
        // see also https://www.pdf-insecurity.org/
        RandomAccessReadBufferedFile raFile = new RandomAccessReadBufferedFile(outFile);
        // If your files are not too large, you can also download the PDF into a byte array
        // with IOUtils.toByteArray() and pass a RandomAccessBuffer() object to the
        // PDFParser constructor.
        PDFParser parser = new PDFParser(raFile, null); // no password!
        try (PDDocument document = parser.parse(false)) {
            for (PDSignature sig : document.getSignatureDictionaries()) {
                COSDictionary sigDict = sig.getCOSObject();
                byte[] contents = sig.getContents();

                // download the signed content
                // we're doing this as a stream, to be able to handle huge files
                try (FileInputStream fis = new FileInputStream(outFile); InputStream signedContentAsStream = new COSFilterInputStream(fis, sig.getByteRange())) {
                    System.out.println("Signature found");
                    if (sig.getName() != null) {
                        System.out.println("Name:     " + sig.getName());
                    }

                    if (sig.getSignDate() != null) {
                        System.out.println("Modified: " + sdf.format(sig.getSignDate().getTime()));
                    }

                    byte[] signedContent = sig.getSignedContent(fis);
                    String subFilter = sig.getSubFilter();

                    if (subFilter != null) {
                        switch (subFilter) {
                            case "adbe.pkcs7.detached":
//                                CMSSignedData signedData = new CMSSignedData(new CMSProcessableByteArray(signedContent), contents);
//                                Store certificatesStore = signedData.getCertificates();
//                                SignerInformationStore signers = signedData.getSignerInfos();
//                                Collection<SignerInformation> c = signers.getSigners();
//                                SignerInformation signerInformation = c.iterator().next();
//
//                                Collection matches = certificatesStore.getMatches(signerInformation.getSID());
//                                X509CertificateHolder certificateHolder = (X509CertificateHolder) matches.iterator().next();
//                                X509Certificate certFromSignedData = new JcaX509CertificateConverter().getCertificate(certificateHolder);
//
//                                if (signerInformation.verify(new JcaSimpleSignerInfoVerifierBuilder().build(certFromSignedData))) {
//                                    System.out.println("Signature verified");
//                                } else {
//                                    System.out.println("Signature verification failed");
//                                }
                                verifyPKCS7(signedContent, contents, sig);
                                break;
                            case "ETSI.CAdES.detached":
                            case "adbe.pkcs7.sha1":
                            case "adbe.x509.rsa_sha1":
                            case "ETSI.RFC3161":
                            default:
                                System.err.println("Unknown certificate type: " + subFilter);
                                break;
                        }
                    } else {
                        throw new IOException("Missing subfilter for cert dictionary");
                    }

                    int[] byteRange = sig.getByteRange();
                    if (byteRange.length != 4) {
                        System.err.println("Signature byteRange must have 4 items");
                    } else {

                    }
                } catch (CMSException e) {
                    throw new RuntimeException(e);
                } catch (GeneralSecurityException | TSPException | OperatorCreationException |
                         CertificateVerificationException e) {
                    throw new RuntimeException(e);
                }
            }
        }
    }

    private void verifyPKCS7(byte[] signedContent, byte[] contents, PDSignature sig) throws CMSException, IOException, GeneralSecurityException, TSPException, OperatorCreationException, CertificateVerificationException {
        // inspiration:
        // http://stackoverflow.com/a/26702631/535646
        // http://stackoverflow.com/a/9261365/535646
        CMSSignedData signedData = new CMSSignedData(new CMSProcessableByteArray(signedContent), contents);
        Store<X509CertificateHolder> certificatesStore = signedData.getCertificates();
        if (certificatesStore.getMatches(null).isEmpty()) {
            throw new IOException("No certificates in signature");
        }
        Collection<SignerInformation> signers = signedData.getSignerInfos().getSigners();
        if (signers.isEmpty()) {
            throw new IOException("No signers in signature");
        }
        SignerInformation signerInformation = signers.iterator().next();
        @SuppressWarnings("unchecked")
        Collection<X509CertificateHolder> matches =
                certificatesStore.getMatches(signerInformation.getSID());
        if (matches.isEmpty()) {
            throw new IOException("Signer '" + signerInformation.getSID().getIssuer() +
                    ", serial# " + signerInformation.getSID().getSerialNumber() +
                    " does not match any certificates");
        }

        X509CertificateHolder certificateHolder = matches.iterator().next();
        X509Certificate certFromSignedData = new JcaX509CertificateConverter().getCertificate(certificateHolder);
        System.out.println("certFromSignedData: " + certFromSignedData);

        SigUtils.checkCertificateUsage(certFromSignedData);

        // Embedded timestamp
        TimeStampToken timeStampToken = SigUtils.extractTimeStampTokenFromSignerInformation(signerInformation);
        if (timeStampToken != null) {
            // tested with QV_RCA1_RCA3_CPCPS_V4_11.pdf
            // https://www.quovadisglobal.com/~/media/Files/Repository/QV_RCA1_RCA3_CPCPS_V4_11.ashx
            // also 021496.pdf and 036351.pdf from digitalcorpora
            SigUtils.validateTimestampToken(timeStampToken);
            X509Certificate certFromTimeStamp = SigUtils.getCertificateFromTimeStampToken(timeStampToken);
            // merge both stores using a set to remove duplicates
            HashSet<X509CertificateHolder> certificateHolderSet = new HashSet<>();
            certificateHolderSet.addAll(certificatesStore.getMatches(null));
            certificateHolderSet.addAll(timeStampToken.getCertificates().getMatches(null));
            SigUtils.verifyCertificateChain(new CollectionStore<>(certificateHolderSet),
                    certFromTimeStamp,
                    timeStampToken.getTimeStampInfo().getGenTime());
            SigUtils.checkTimeStampCertificateUsage(certFromTimeStamp);

            // compare the hash of the signature with the hash in the timestamp
            byte[] tsMessageImprintDigest = timeStampToken.getTimeStampInfo().getMessageImprintDigest();
            String hashAlgorithm = timeStampToken.getTimeStampInfo().getMessageImprintAlgOID().getId();
            byte[] sigMessageImprintDigest = MessageDigest.getInstance(hashAlgorithm).digest(signerInformation.getSignature());
            if (Arrays.equals(tsMessageImprintDigest, sigMessageImprintDigest)) {
                System.out.println("timestamp signature verified");
            } else {
                System.err.println("timestamp signature verification failed");
            }
        }

        try {
            if (sig.getSignDate() != null) {
                certFromSignedData.checkValidity(sig.getSignDate().getTime());
                System.out.println("Certificate valid at signing time");
            } else {
                System.err.println("Certificate cannot be verified without signing time");
            }
        } catch (CertificateExpiredException ex) {
            System.err.println("Certificate expired at signing time");
        } catch (CertificateNotYetValidException ex) {
            System.err.println("Certificate not yet valid at signing time");
        }

        // usually not available
        if (signerInformation.getSignedAttributes() != null) {
            // From SignedMailValidator.getSignatureTime()
            Attribute signingTime = signerInformation.getSignedAttributes().get(CMSAttributes.signingTime);
            if (signingTime != null) {
                Time timeInstance = Time.getInstance(signingTime.getAttrValues().getObjectAt(0));
                try {
                    certFromSignedData.checkValidity(timeInstance.getDate());
                    System.out.println("Certificate valid at signing time: " + timeInstance.getDate());
                } catch (CertificateExpiredException ex) {
                    System.err.println("Certificate expired at signing time");
                } catch (CertificateNotYetValidException ex) {
                    System.err.println("Certificate not yet valid at signing time");
                }
            }
        }

        if (signerInformation.verify(new JcaSimpleSignerInfoVerifierBuilder().
                setProvider(SecurityProvider.getProvider()).build(certFromSignedData))) {
            System.out.println("Signature verified");
        } else {
            System.out.println("Signature verification failed");
        }

        if (CertificateVerifier.isSelfSigned(certFromSignedData)) {
            System.err.println("Certificate is self-signed, LOL!");
        } else {
            System.out.println("Certificate is not self-signed");

            if (sig.getSignDate() != null) {
                SigUtils.verifyCertificateChain(certificatesStore, certFromSignedData, sig.getSignDate().getTime());
            } else {
                System.err.println("Certificate cannot be verified without signing time");
            }
        }
    }
}
