# 1. Use OpenSSL:

- A certificate authority (CA) is responsible for storing, signing, and issuing digital certificates

- This means we will first generate a private key and a self-signed certificate for the certificate authority (CA):

  openssl req -x509 \<!-- Public key certificate format -->
  -sha256 \<!-- Uses a sha-256 message digest algorithm -->
  -newkey rsa:4096 \<!-- Generates a private key and its certificate request -->
  -days 365 \<!-- Expiration date -->
  -keyout ca-key.pem \<!-- Private key file -->
  -out ca-cert.pem \<!-- Certificate file -->
  -subj "/C=TR/ST=EURASIA/L=ISTANBUL/O=Software/OU=Microservices/CN=\*.microservices.dev/emailAddress=huseyin@microservices.dev"<!-- Adds identity information to certificate -->
  -nodes<!-- "no DES" means the private key will not be encrypted -->

- The -subj parameter contains identity information about the certificate:

  /C is used for country
  /ST is the state information
  /L states city information
  /O means organization
  /OU is for the organization unit to explain which department
  /CN is used for the domain name, the short version of common name
  /emailAddress is used for an email address to contact the certificate owner

# 2. Verify the generated self-certificate:

- You can verify the generated self-certificate for the CA with the following command:

  openssl x509 -in ca-cert.pem -noout -text

# 3. Private key and certificate signing request:

- Once you verify it, we can proceed with the private key and certificate signing request:

  openssl req \<!-- Certificate signing request -->
  -newkey rsa:4096 \
   -keyout server-key.pem \<!-- The location of the private key -->
  -out server-req.pem \<!-- The location of the certificate request -->
  -subj "/C=TR/ST=EURASIA/L=ISTANBUL/O=Microservices/OU=PaymentService/CN=\*.microservices.dev/emailAddress=huseyin@microservices.dev" \
   -nodes \
   -sha256

- Then we will use CA’s private key to sign the request:
  openssl x509 \
   -req -in server-req.pem \<!-- Passes the sign request -->
  -days 60 \
   -CA ca-cert.pem \<!-- CA's certificate -->
  -CAkey ca-key.pem \<!-- CA's private key -->
  -CAcreateserial \<!-- Generates the next serial number for the certificate -->
  -out server-cert.pem \
   -extfile server-ext.cnf \<!-- Additional configs for the certificate -->
  -sha256

- An example configuration for ext file option is as follows:

  subjectAltName=DNS:\*.microservices.dev,DNS:\*.microservices.dev,IP:0.0.0.0

- Now you can verify the server’s self-signed certificate:

  openssl x509 -in server-cert.pem -noout -text

# 4. mTLS:

- For mTLS communication, we need to generate a certificate signing request for the client side, so let’s generate a private key and this self-signed certificate:

  openssl req \
   -newkey rsa:4096 \
   -keyout client-key.pem \
   -out client-req.pem \
   -subj "/C=TR/ST=EURASIA/L=Istanbul/O=Microservices/OU=OrderService/CN=\*.microservices.dev/emailAddress=huseyin@microservices.dev" \
   -nodes \
   -sha256

- Now, let’s sign it using the CA’s private key:

  openssl x509 \
   -req -in client-req.pem \
   -sha256 \
   -days 60 \
   -CA ca.crt \
   -CAkey ca.key \
   -CAcreateserial \
   -out client.crt \
   -extfile client-ext.cnf

- Finally, you can verify the client certificate with the following command:

  openssl x509 -in client-cert.pem -noout -text
