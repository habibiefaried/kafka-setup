#!/bin/bash

set -o nounset \
    -o errexit

printf "Deleting previous (if any)..."
rm -rf secrets
mkdir secrets
mkdir -p tmp
echo " OK!"
# Generate CA key
printf "Creating CA..."
openssl req -new -x509 -keyout tmp/datahub-ca.key -out tmp/datahub-ca.crt -days 365 -subj '/CN=ca.datahub/OU=test/O=datahub/L=paris/C=fr' -passin pass:datahub -passout pass:datahub >/dev/null 2>&1

echo " OK!"

for i in 'broker' 'producer' 'consumer' 'schema-registry'; do
    printf "Creating cert and keystore of $i..."
    # Create keystores
    keytool -genkey -noprompt \
        -alias $i \
        -dname "CN=$i, OU=test, O=datahub, L=paris, C=fr" \
        -keystore secrets/$i.keystore.jks \
        -keyalg RSA \
        -storepass datahub \
        -keypass datahub >/dev/null 2>&1

    # Create CSR, sign the key and import back into keystore
    keytool -keystore secrets/$i.keystore.jks -alias $i -certreq -file tmp/$i.csr -storepass datahub -keypass datahub >/dev/null 2>&1

    openssl x509 -req -CA tmp/datahub-ca.crt -CAkey tmp/datahub-ca.key -in tmp/$i.csr -out tmp/$i-ca-signed.crt -days 365 -CAcreateserial -passin pass:datahub >/dev/null 2>&1

    keytool -keystore secrets/$i.keystore.jks -alias CARoot -import -noprompt -file tmp/datahub-ca.crt -storepass datahub -keypass datahub >/dev/null 2>&1

    keytool -keystore secrets/$i.keystore.jks -alias $i -import -file tmp/$i-ca-signed.crt -storepass datahub -keypass datahub >/dev/null 2>&1

    # Create truststore and import the CA cert.
    keytool -keystore secrets/$i.truststore.jks -alias CARoot -import -noprompt -file tmp/datahub-ca.crt -storepass datahub -keypass datahub >/dev/null 2>&1
    echo " OK!"
done

echo "datahub" >secrets/cert_creds
cp tmp/datahub-ca.crt kafka_server_cert.pem
rm -rf tmp

echo "SUCCEEDED"

# Export client keystore to PKCS12
echo "Exporting client keystore to PKCS12..."
keytool -importkeystore -srckeystore secrets/consumer.keystore.jks -destkeystore consumer.pkcs12 -deststoretype PKCS12 -srcalias consumer -deststorepass datahub -destkeypass datahub -srcstorepass datahub

# Extract client certificate and key
echo "Extracting client certificate and key..."
openssl pkcs12 -in consumer.pkcs12 -nokeys -out "kafka_client_cert.pem" -passin pass:datahub
openssl pkcs12 -in consumer.pkcs12 -nocerts -nodes -out "kafka_client_key.pem" -passin pass:datahub

sleep 5

docker-compose up -d
