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
//        pdfdemo.createImageButton();
//        pdfdemo.updateImageButton();
        pdfdemo.insertImage("output.pdf");
    }
}
