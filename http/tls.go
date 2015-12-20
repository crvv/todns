package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/crvv/todns/config"
	"io/ioutil"
	"log"
	"net"
)

func newTlsListener() net.Listener {
	certificate, err := tls.LoadX509KeyPair(config.GetCrt(), config.GetKey())
	if err != nil {
		log.Fatalln(err)
	}
	clientCrt, err := ioutil.ReadFile(config.GetClientCrt())
	if err != nil {
		log.Fatalln(err)
	}
	block, _ := pem.Decode(clientCrt)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalln(err)
	}
	clientCAs := x509.NewCertPool()
	clientCAs.AddCert(cert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		NextProtos:   []string{"http"},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	}
	listener, err := tls.Listen("tcp", config.GetHttpsListen(), tlsConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return listener
}
