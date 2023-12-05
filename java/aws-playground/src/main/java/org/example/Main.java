package org.example;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProviderChain;
import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;
import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.ListObjectVersionsResponse;
import software.amazon.awssdk.services.s3.model.ObjectVersion;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;
import java.util.Collections;
import java.util.UUID;

public class Main {

    // You must create aws profile first
    //    aws configure --profile localstack
    // Alternative,
    //    StaticCredentialsProvider.create(AwsBasicCredentials.create(ACCESS_KEY, SECRET_KEY))
    public static final AwsCredentialsProviderChain CREDENTIALS_PROVIDER_CHAIN = AwsCredentialsProviderChain.of(ProfileCredentialsProvider.builder().profileName("localstack").build(), DefaultCredentialsProvider.create());
    public static final URI ENDPOINT = URI.create("http://s3.localhost.localstack.cloud:4566");

    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    public static void main(String[] args) {
//        String bucketName = "b-" + UUID.randomUUID();
//        String keyName = "k-" + UUID.randomUUID();
        String bucketName = "test-bucket";
        String keyName = "test-key";
        File file = new File("input.pdf");

        try (S3Client s3Client = S3Client.builder().region(Region.US_EAST_1).
                credentialsProvider(CREDENTIALS_PROVIDER_CHAIN).
                endpointOverride(ENDPOINT).
                build()) {

            S3Presigner presigner = S3Presigner.builder().region(Region.US_EAST_1).
                    credentialsProvider(CREDENTIALS_PROVIDER_CHAIN).
                    endpointOverride(ENDPOINT).build();
            S3Cli cli = new S3Cli(s3Client, presigner);

            cli.createBucket(bucketName, true);

            try {
                URL presignedUrl = cli.createPresignedUrl(bucketName, keyName, Files.probeContentType(file.toPath()), Collections.EMPTY_MAP);
                cli.useHttpClientToPut(presignedUrl, file, Files.probeContentType(file.toPath()), Collections.EMPTY_MAP);
                cli.useHttpUrlConnectionToPut(presignedUrl, file, Files.probeContentType(file.toPath()), Collections.EMPTY_MAP);

                // Get it
                cli.getPresignedUrl(bucketName, keyName, "/tmp/input_downloaded.pdf");

                // List object versions
                ListObjectVersionsResponse versions = cli.listObjectVersions(bucketName);
                for (ObjectVersion version : versions.versions()) {
                    logger.info(version.toString());
                }

                logger.info("-----------------------------------------------------------");

                // Get the last item -> the oldest object
                logger.info(versions.versions().getLast().toString());
            } catch (IOException e) {
                throw new RuntimeException(e);
            } finally {
//                cli.deleteObject(bucketName, keyName);
//                cli.deleteBucket(bucketName);
            }
        }
    }
}