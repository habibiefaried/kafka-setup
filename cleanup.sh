#!/bin/bash
docker-compose down -v
rm -rf ./secrets
rm -rf *.pkcs12
rm -rf *.pem
