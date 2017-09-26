#!/bin/sh

rm -rf /etc/certs
signdomain=HOST
mkdir -p /etc/certs
cat << EOF > /etc/certs/extfile.cnf
[ v3_ca ]
subjectAltName = IP:$signdomain
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = CA:true
EOF
openssl req -nodes -subj "/C=CN/ST=BeiJing/L=BeiJing/CN=$signdomain" -newkey rsa:4096 -keyout /etc/certs/$signdomain.key -out /etc/certs/$signdomain.csr
openssl x509 -req -days 3650 -in /etc/certs/$signdomain.csr -signkey /etc/certs/$signdomain.key -out /etc/certs/$signdomain.crt -extfile /etc/certs/extfile.cnf -extensions v3_ca
