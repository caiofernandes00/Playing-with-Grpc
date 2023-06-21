# Tutorial -> https://dev.to/techschoolguru/how-to-create-sign-ssl-tls-certificates-2aai
rm *.pem
rm *.srl

###################### CA ######################
# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -nodes -keyout ca-key.pem -out ca-cert.pem -days 365 -subj "/C=BR/ST=SP/L=São Paulo/O=Authority/OU=AuthorityOffice/CN=*.homeoffice.org/emailAddress=certificate@authorityoffice.org"
openssl x509 -in ca-cert.pem -noout -text


###################### Server ######################
# Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -nodes -subj "/C=BR/ST=SP/L=São Paulo/O=HomeOffice/OU=HomeOffice/CN=*.myhome.com/emailAddress=caiow.wk2@pm.me"

# Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
openssl x509 -in server-cert.pem -noout -text

###################### Client ######################
# Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -nodes -subj "/C=BR/ST=SP/L=São Paulo/O=HomeOffice/OU=HomeOffice/CN=*.myhome.com/emailAddress=caiow.wk2@pm.me"

# Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf
openssl x509 -in client-cert.pem -noout -text