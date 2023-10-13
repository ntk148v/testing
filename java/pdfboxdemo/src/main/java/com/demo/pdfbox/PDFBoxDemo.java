package com.demo.pdfbox;

import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

import javax.security.auth.x500.X500Principal;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.spec.RSAKeyGenParameterSpec;
import java.util.Date;

public class PDFBoxDemo {
    private static final int KEY_SIZE = 2048;
    private final KeyPair keyPair;
    private X509Certificate certificate;

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

    private void generateCertificate() {
        // Add Bouncy Castle as a Security Provider
        Security.addProvider(new BouncyCastleProvider());
        X509v3CertificateBuilder certBuilder = new JcaX509v3CertificateBuilder(
                new X500Name("CN=Test"),
                BigInteger.valueOf(1), // serial
                new Date(System.currentTimeMillis() - 1000L * 60 * 60 * 24 * 30),
                new Date(System.currentTimeMillis() + (1000L * 60 * 60 * 24 * 30)),
                new X500Name("CN=Test"),
                this.keyPair.getPublic()
        );

        ContentSigner contentSigner;
        try {
            contentSigner = new JcaContentSignerBuilder("SHA256WithRSAEncryption").
                    build(this.keyPair.getPrivate());
            this.certificate = new JcaX509CertificateConverter().setProvider("BC")
                    .getCertificate(certBuilder.build(contentSigner));
        } catch (OperatorCreationException e) {
            throw new RuntimeException(e);
        } catch (CertificateException e) {
            throw new RuntimeException(e)
        }
    }

    public static void main(String[] args) {

    }
}
