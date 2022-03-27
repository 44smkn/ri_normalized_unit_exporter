# AWS Reserved Instance Exporter

![test](https://github.com/44smkn/aws_ri_exporter/actions/workflows/test.yaml/badge.svg)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Metrics exporter for aws reserved instance normalized unit per hour.

## Installation and Usage

The `aws_ri_exporter` listens on HTTP port **9981** by default. See the `--help` output for more options.

You will need to have AWS API credentials configured. What works for AWS CLI, should be sufficient. You can use [~/.aws/credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) or [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-set).

### Setup IAM role

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "rds:DescribeDBInstances",
        "rds:DescribeReservedDBInstances"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

### Docker

```sh
docker run -d -e AWS_REGION=ap-northeast-1 \ 
  ghcr.io/44smkn/ri-normalized-unit-exporter:latest --log.level=debug
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aws-ri-exporter
  labels:
    app: aws-ri-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: aws-ri-exporter
  template:
    metadata:
      labels:
        app: aws-ri-exporter
    spec:
      containers:
      - name: aws-ri-exporter
        image: ghcr.io/44smkn/zenhub_exporter:latest
        ports:
        - containerPort: 9981
        env:
        - name: AWS_REGION
          value: ap-northeast-1
```

## Collectors

Collectors are enabled by providing a `--collector.<name>` flag.
Collectors that are enabled by default can be disabled by providing a `--no-collector.<name>` flag.

| Collector | Metrics                                         | Description                                    |
|-----------|-------------------------------------------------|------------------------------------------------|
| rds       | `aws_ri_rds_running_instance_normalized_unit`   | Normalized Units for each running RDS instance |
| rds       | `aws_ri_rds_active_reservation_normalized_unit` | Normalized Units for each active reservation   |
