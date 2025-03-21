package com.demo.pdfbox;

import java.io.IOException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.UnrecoverableKeyException;
import java.security.cert.CertificateException;

public class Main {
    public static void main(String[] args) throws UnrecoverableKeyException, CertificateException, IOException, KeyStoreException, NoSuchAlgorithmException {
        PDFBoxDemo pdfdemo = new PDFBoxDemo();
//        pdfdemo.sign1();
//        pdfdemo.verify("output1.pdf");

//        pdfdemo.sign2();
//        pdfdemo.verify("output2.pdf");
        pdfdemo.createImageButton();
//        pdfdemo.fillImageButton("/tmp/output.pdf", 1);
//        pdfdemo.updateImageButton("/tmp/input.pdf", 2, 26, 9, 100, 100);
//        pdfdemo.insertImage("/tmp/output.pdf");
//        pdfdemo.updateImageButton("/tmp/input-create-form-signed-pades-baseline-b.pdf");

//        pdfdemo.genThumbnail("jpg", 0f, 0);
    }
}
