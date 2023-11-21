package org.example;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.model.S3Exception;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;
import software.amazon.awssdk.services.s3.presigner.model.PresignedPutObjectRequest;
import software.amazon.awssdk.services.s3.presigner.model.PutObjectPresignRequest;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URISyntaxException;
import java.net.URL;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;
import java.nio.file.Path;
import java.time.Duration;
import java.util.Map;
import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.time.Duration;

import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.model.GetObjectRequest;
import software.amazon.awssdk.services.s3.model.S3Exception;
import software.amazon.awssdk.services.s3.presigner.model.GetObjectPresignRequest;
import software.amazon.awssdk.services.s3.presigner.model.PresignedGetObjectRequest;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;
import software.amazon.awssdk.utils.IoUtils;


public class S3Cli {
    public S3Cli(S3Presigner presigner) {
        this.presigner = presigner;
    }

    private final S3Presigner presigner;

    private final static Logger logger = LoggerFactory.getLogger(S3Cli.class);

    /**
     * Create a presigned URL for uploading with a PUT request.
     *
     * @param bucketName  - The name of the bucket.
     * @param keyName     - The name of the object.
     * @param contentType - The content type of the object.
     * @param metadata    - The metadata to store with the object.
     * @return - The presigned URL for an HTTP PUT.
     */
    public URL createPresignedUrl(String bucketName, String keyName, String contentType, Map<String, String> metadata) {
        try {

            PutObjectRequest objectRequest = PutObjectRequest.builder()
                    .bucket(bucketName)
                    .key(keyName)
                    .contentType(contentType)
                    .metadata(metadata)
                    .build();

            PutObjectPresignRequest presignRequest = PutObjectPresignRequest.builder()
                    .signatureDuration(Duration.ofMinutes(10))  // The URL will expire in 10 minutes.
                    .putObjectRequest(objectRequest)
                    .build();

            PresignedPutObjectRequest presignedRequest = presigner.presignPutObject(presignRequest);
            String myURL = presignedRequest.url().toString();
            logger.info("Presigned URL to upload a file to: [{}]", myURL);
            logger.info("Which HTTP method needs to be used when uploading a file: [{}]", presignedRequest.httpRequest().method());

            return presignedRequest.url();
        } catch (Exception e) {
            logger.error(e.getMessage(), e);
            return null;
        }
    }

    /**
     * Use the JDK HttpURLConnection (since v1.1) class to do the upload, but you can
     * use any HTTP client.
     *
     * @param presignedUrl - The presigned URL.
     * @param fileToPut    - The file to upload.
     * @param contentType  - The content type of the file.
     * @param metadata     - The metadata to store with the object.
     */
    public void useHttpUrlConnectionToPut(URL presignedUrl, File fileToPut, String contentType, Map<String, String> metadata) {
        logger.info("Begin [{}] upload", contentType);
        try {
            HttpURLConnection connection = (HttpURLConnection) presignedUrl.openConnection();
            connection.setDoOutput(true);
            connection.setRequestProperty("Content-Type", contentType);
            metadata.forEach((k, v) -> connection.setRequestProperty("x-amz-meta-" + k, v));
            connection.setRequestMethod("PUT");
            OutputStream out = connection.getOutputStream();

            try (RandomAccessFile file = new RandomAccessFile(fileToPut, "r");
                 FileChannel inChannel = file.getChannel()) {
                ByteBuffer buffer = ByteBuffer.allocate(8192); //Buffer size is 8k

                while (inChannel.read(buffer) > 0) {
                    buffer.flip();
                    for (int i = 0; i < buffer.limit(); i++) {
                        out.write(buffer.get());
                    }
                    buffer.clear();
                }
            } catch (IOException e) {
                logger.error(e.getMessage(), e);
            }

            out.close();
            connection.getResponseCode();
            logger.info("HTTP response code is " + connection.getResponseCode());

        } catch (S3Exception | IOException e) {
            logger.error(e.getMessage(), e);
        }
    }

    /**
     * Use the JDK HttpClient (since v11) class to do the upload, but you can
     * use any HTTP client.
     *
     * @param presignedUrl - The presigned URL.
     * @param fileToPut    - The file to upload.
     * @param contentType  - The content type of the file.
     * @param metadata     - The metadata to store with the object.
     */
    public void useHttpClientToPut(URL presignedUrl, File fileToPut, String contentType, Map<String, String> metadata) {
        logger.info("Begin [{}] upload", contentType);

        HttpRequest.Builder requestBuilder = HttpRequest.newBuilder();
        metadata.forEach((k, v) -> requestBuilder.header("x-amz-meta-" + k, v));

        HttpClient httpClient = HttpClient.newHttpClient();
        try {
            final HttpResponse<Void> response = httpClient.send(requestBuilder
                            .uri(presignedUrl.toURI())
                            .header("Content-Type", contentType)
                            .PUT(HttpRequest.BodyPublishers.ofFile(Path.of(fileToPut.toURI())))
                            .build(),
                    HttpResponse.BodyHandlers.discarding());

            logger.info("HTTP response code is " + response.statusCode());

        } catch (URISyntaxException | InterruptedException | IOException e) {
            logger.error(e.getMessage(), e);
        }
    }

    public void getPresignedUrl(String bucketName, String keyName, String output) {

        try {
            GetObjectRequest getObjectRequest = GetObjectRequest.builder()
                    .bucket(bucketName)
                    .key(keyName)
                    .build();

            GetObjectPresignRequest getObjectPresignRequest = GetObjectPresignRequest.builder()
                    .signatureDuration(Duration.ofMinutes(60))
                    .getObjectRequest(getObjectRequest)
                    .build();

            PresignedGetObjectRequest presignedGetObjectRequest = presigner.presignGetObject(getObjectPresignRequest);
            String theUrl = presignedGetObjectRequest.url().toString();
            logger.info("Presigned URL: " + theUrl);
            HttpURLConnection connection = (HttpURLConnection) presignedGetObjectRequest.url().openConnection();
            presignedGetObjectRequest.httpRequest().headers().forEach((header, values) -> {
                values.forEach(value -> {
                    connection.addRequestProperty(header, value);
                });
            });

            // Send any request payload that the service needs (not needed when isBrowserExecutable is true).
            if (presignedGetObjectRequest.signedPayload().isPresent()) {
                connection.setDoOutput(true);

                try (InputStream signedPayload = presignedGetObjectRequest.signedPayload().get().asInputStream();
                     OutputStream httpOutputStream = connection.getOutputStream()) {
                    IoUtils.copy(signedPayload, httpOutputStream);
                }
            }

            // Download the result of executing the request.
            try (InputStream content = connection.getInputStream()) {
                logger.info("Service returned response: ");
                if (output != null) {
                    IoUtils.copy(content, new FileOutputStream(output));
                } else {
                    IoUtils.copy(content, System.out);
                }
            }

        } catch (S3Exception | IOException e) {
            logger.error(e.getMessage(), e);
        }
    }
}
