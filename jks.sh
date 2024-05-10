#!/bin/bash

# Set variables
CA_ALIAS="kafkaCA"
CA_KEYSTORE="kafkaCA.keystore.jks"
SERVER_ALIAS="kafkaServer"
KEYSTORE="kafka.keystore.jks"
TRUSTSTORE="kafka.truststore.jks"
PASSWORD="Ch4ng3It"
DNAME_CA="CN=kafkaCA, OU=IT, O=YourCompany, L=YourCity, ST=YourState, C=US"
DNAME_SERVER="CN=kafkaServer, OU=IT, O=YourCompany, L=YourCity, ST=YourState, C=US"
VALIDITY=365
KEY_SIZE=2048

# Create CA keystore with the CA keypair
echo "Creating CA keystore and key pair..."
keytool -genkeypair -alias $CA_ALIAS -keyalg RSA -keysize $KEY_SIZE -keystore $CA_KEYSTORE -validity $VALIDITY -storepass $PASSWORD -keypass $PASSWORD -dname "$DNAME_CA"

# Generate the CA certificate
echo "Generating CA certificate..."
keytool -export -alias $CA_ALIAS -keystore $CA_KEYSTORE -file $CA_ALIAS.cer -storepass $PASSWORD

# Create server keystore with the keypair
echo "Creating server keystore and key pair..."
keytool -genkeypair -alias $SERVER_ALIAS -keyalg RSA -keysize $KEY_SIZE -keystore $KEYSTORE -validity $VALIDITY -storepass $PASSWORD -keypass $PASSWORD -dname "$DNAME_SERVER"

# Create a certificate signing request (CSR) for the server
echo "Creating CSR for the server..."
keytool -certreq -alias $SERVER_ALIAS -keystore $KEYSTORE -file $SERVER_ALIAS.csr -storepass $PASSWORD

# Sign the server CSR with the CA key
echo "Signing the server certificate with the CA..."
keytool -gencert -alias $CA_ALIAS -keystore $CA_KEYSTORE -infile $SERVER_ALIAS.csr -outfile $SERVER_ALIAS.cer -validity $VALIDITY -storepass $PASSWORD

# Import the CA certificate to the server keystore
echo "Importing CA certificate to server keystore..."
keytool -importcert -alias $CA_ALIAS -keystore $KEYSTORE -file $CA_ALIAS.cer -storepass $PASSWORD -noprompt

# Import the signed server certificate to the server keystore
echo "Importing signed server certificate to server keystore..."
keytool -importcert -alias $SERVER_ALIAS -keystore $KEYSTORE -file $SERVER_ALIAS.cer -storepass $PASSWORD -noprompt

# Create the truststore and import the CA certificate
echo "Creating truststore and importing the CA certificate..."
keytool -importcert -alias $CA_ALIAS -file $CA_ALIAS.cer -keystore $TRUSTSTORE -storepass $PASSWORD -noprompt

echo "JKS creation and CA setup completed."

#### Client cert generation

keytool -genkeypair -alias kafkaClient -keyalg RSA -keystore kafka.client.keystore.jks -storepass changeit -validity 365 -keysize 2048 -dname "CN=KafkaClient, OU=IT, O=YourCompany, L=City, S=State, C=US"
