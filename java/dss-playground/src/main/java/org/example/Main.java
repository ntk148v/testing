package org.example;

 import eu.europa.esig.dss.detailedreport.DetailedReport;
 import eu.europa.esig.dss.diagnostic.DiagnosticData;
 import eu.europa.esig.dss.enumerations.DigestAlgorithm;
 import eu.europa.esig.dss.model.SignatureValue;
 import eu.europa.esig.dss.model.ToBeSigned;
 import eu.europa.esig.dss.model.x509.CertificateToken;
 import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
 import eu.europa.esig.dss.token.MSCAPISignatureToken;
 import eu.europa.esig.dss.utils.Utils;
 import eu.europa.esig.dss.spi.CertificateExtensionsUtils;
 import eu.europa.esig.dss.validation.CertificateValidator;
 import eu.europa.esig.dss.validation.CertificateVerifier;
 import eu.europa.esig.dss.validation.CommonCertificateVerifier;
 import eu.europa.esig.dss.service.crl.OnlineCRLSource;
 import eu.europa.esig.dss.service.ocsp.OnlineOCSPSource;
 import eu.europa.esig.dss.validation.reports.CertificateReports;
 import eu.europa.esig.dss.simplecertificatereport.SimpleCertificateReport;

 import java.util.List;


public class Main {
    public static void main(String[] args) {
        try (MSCAPISignatureToken token = new MSCAPISignatureToken()) {

            List<DSSPrivateKeyEntry> keys = token.getKeys();
            for (DSSPrivateKeyEntry entry : keys) {

                CertificateToken certificateToken = entry.getCertificate();
                CertificateVerifier cv = new CommonCertificateVerifier();
                cv.setOcspSource(new OnlineOCSPSource());
                cv.setCrlSource(new OnlineCRLSource());
                // We create an instance of the CertificateValidator with the certificate
                CertificateValidator validator = CertificateValidator.fromCertificate(certificateToken);
                validator.setCertificateVerifier(cv);
                CertificateReports certificateReports = validator.validate();
                // We have 3 reports
                // The diagnostic data which contains all used and static data
                DiagnosticData diagnosticData = certificateReports.getDiagnosticData();

                // The detailed report which is the result of the process of the diagnostic data and the validation policy
                DetailedReport detailedReport = certificateReports.getDetailedReport();

                // The simple report is a summary of the detailed report or diagnostic data (more user-friendly)
                SimpleCertificateReport simpleReport = certificateReports.getSimpleReport();
                System.out.println(simpleReport);
            }

//            ToBeSigned toBeSigned = new eu.europa.esig.dss.model.ToBeSigned("Hello world".getBytes());
//            SignatureValue signatureValue = token.sign(toBeSigned, DigestAlgorithm.SHA256, keys.get(0));
//
//            System.out.println("Signature value : " + Utils.toBase64(signatureValue.getValue()));
        }
    }
}