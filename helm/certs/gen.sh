#!/bin/sh
openssl req -config ./config -new -x509 -sha256 -newkey rsa:2048 -nodes -keyout envoy.key.pem -days 365 -out envoy.cert.pem
