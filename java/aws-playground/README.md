# AWS Playground

- Run [localstack](https://www.localstack.cloud/) container:

```shell
 docker run --rm --name local-s3 -p 4566:4566 -e SERVICES=s3 -e EAGER_SERVICE_LOADING=1 -e START_WEB=0 -d localstack/localstack
```

- Configure localstack profile:

```shell
aws configure --profile localstack
AWS Access Key ID [None]: test-key
AWS Secret Access Key [None]: test-secret
Default region name [None]: us-east-1
Default output format [None]:
```
