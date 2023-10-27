package org.example;

import eu.europa.esig.dss.model.SignatureValue;
import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
import eu.europa.esig.dss.token.Pkcs11SignatureToken;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;


import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import java.io.IOException;
import java.math.BigInteger;
import java.security.*;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.spec.RSAKeyGenParameterSpec;
import java.util.Date;
import java.util.List;

public class Main {
    private static String HSM_PIN = "1234";
    private static final int KEY_SIZE = 2048;
    private static String KEY_PASSWD = "1234";
    private static String KEY_ALIAS = "keycuakien";

    public static void main(String[] args) {
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

        try (Pkcs11SignatureToken token = new Pkcs11SignatureToken("/usr/local/lib/softhsm/libsofthsm2.so",
                new KeyStore.PasswordProtection("1234".toCharArray()), 1651866022)) {
            List<DSSPrivateKeyEntry> keys = token.getKeys();
            for (DSSPrivateKeyEntry entry : keys) {
                System.out.println(entry.getCertificate().getCertificate());
                System.out.println(entry.getCertificate().isSelfSigned());
            }
        }
    }


}