# kafka-setup
Kafka setup with JKS with golang testcases

## Requirements

Install keytool

```
apt install openjdk-19-jre-headless -y
```

## Run

1. To spin up cluster with ssl

```
make build
```

2. To run tests

```
make test
```

3. To cleanup

```
make cleanup
```