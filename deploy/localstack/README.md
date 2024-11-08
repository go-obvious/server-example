# (Bottling AWS) Local Stack Deployment

Localstack is a fully functional local AWS cloud stack. It allows you to develop and test your cloud applications offline, without the need to connect to actual AWS services. This can significantly reduce development costs and provide a faster feedback loop for debugging and iterating on your code.

By using Localstack, you can simulate a wide range of AWS services on your local machine, ensuring that your application behaves as expected before deploying it to the cloud. In this setup, S3 in Localstack will be used for the remotewrite API, enabling you to test and validate your remote write operations locally.

## Quick Start

### 1. Deploy Localstack

Deploy Localstack to your Kubernetes cluster by running:

```sh
deploy.sh
```

**Expected Output:**

```sh
Release "localstack" does not exist. Installing it now.
NAME: localstack
LAST DEPLOYED: Tue Oct 15 18:32:53 2024
NAMESPACE: localstack
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  NOTE: It may take a few minutes for the LoadBalancer IP to be available.
      You can watch the status of by running 'kubectl get --namespace "aws" svc -w localstack'
  export SERVICE_IP=$(kubectl get svc --namespace "aws" localstack --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")
  echo http://$SERVICE_IP:4566
```

### 2. Configure Environment Variables

Set up the necessary environment variables:

```sh
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export AWS_REGION="us-east-1"
export AWS_DEFAULT_REGION="us-east-1"
export AWS_ENDPOINT="http://localhost:4566"
```

### 3. Test Localstack

Run the following commands to test Localstack:

_make an alias for QoL_

```sh
alias awslocal="aws --endpoint-url=$AWS_ENDPOINT"
```

_Try making a bucket, listing it, then delete it._
```sh
awslocal s3api create-bucket --bucket foobar
awslocal s3 ls
awslocal s3api delete-bucket --bucket foobar
```

> Note, `AWS_ENDPOINT` can be used when creating `awsConfig` objects used in aws session creation - which allows us to talk to the localstack s3 implementation.

### 4. Uninstall Localstack

To uninstall Localstack, run:

```sh
delete.sh
```