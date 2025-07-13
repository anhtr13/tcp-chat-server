package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func LoadCerts() (tls.Certificate, error) {
	home_dir, _ := os.UserHomeDir()
	cert_dir := fmt.Sprintf("%s/.tcp-chat/certs", home_dir)

	if _, err := os.Stat(cert_dir); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(cert_dir, os.ModePerm)
	}

	script_file := fmt.Sprintf("%s/cert.sh", cert_dir)
	if _, err := os.Stat(script_file); errors.Is(err, os.ErrNotExist) {
		os.WriteFile(script_file, []byte(`#!/bin/sh

cert_dir=~/.tcp-chat/certs

rm -rf $cert_dir/*.pem
rm -rf $cert_dir/*.srl

echo "CA's self-signed certificate:"
openssl req -x509 -nodes -newkey rsa:4096 -days 3650 -keyout $cert_dir/ca-key.pem -out $cert_dir/ca-cert.pem


echo "Server's self-signed certificate:"
openssl req -nodes -newkey rsa:4096 -keyout $cert_dir/server-key.pem -out $cert_dir/server-req.pem

openssl x509 -req -in $cert_dir/server-req.pem -CA $cert_dir/ca-cert.pem -CAkey $cert_dir/ca-key.pem -CAcreateserial -out $cert_dir/server-cert.pem -days 3650
`), os.ModePerm)
	}

	cert_file := fmt.Sprintf("%s/server-cert.pem", cert_dir)
	key_file := fmt.Sprintf("%s/server-key.pem", cert_dir)

	_, err1 := os.Stat(cert_file)
	_, err2 := os.Stat(cert_file)
	if errors.Is(err1, os.ErrNotExist) || errors.Is(err2, os.ErrNotExist) {
		fmt.Println("No SSL certifictes found, creating certifictes...")
		fmt.Println()
		cmd := exec.Command("/bin/sh", script_file)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return tls.Certificate{}, err
		}
	}

	return tls.LoadX509KeyPair(cert_file, key_file)
}
