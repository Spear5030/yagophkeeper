rm -rf ./cert/*
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ./cert/ca-key.pem -out ./cert/ca-cert.pem -subj "/C=RU/L=Saint-Petersburg/O=Server/OU=Education/CN=*/emailAddress=spear5030@gmail.com"
openssl x509 -in ./cert/ca-cert.pem -noout -text
openssl req -newkey rsa:4096 -nodes -keyout ./cert/server-key.pem -out ./cert/server-req.pem -subj "/C=RU/L=Saint-Petersburg/O=Server/OU=Education/CN=*/emailAddress=spear5030@gmail.com"
openssl x509 -req  -in ./cert/server-req.pem -days 60 -CA ./cert/ca-cert.pem -CAkey ./cert/ca-key.pem -CAcreateserial -out ./cert/server-cert.pem -extfile ./cert/server-ext.cnf
openssl x509 -in ./cert/server-cert.pem -noout -text
