package proxy

import (
	"fmt"
	"log"
	"net"

	"github.com/things-go/go-socks5"
)

var listener net.Listener

type Config struct {
	ListenAddr string
	LocalNames []string
}

func Start(cfg *Config) error {
	l, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return fmt.Errorf("cannot start TCP listener on %s: %w", cfg.ListenAddr, err)
	}
	listener = l // using own listener to be able to run async
	socksServer := socks5.NewServer(socks5.WithResolver(&resolver{
		localNames: cfg.LocalNames,
	}))
	log.Printf("Started proxy server on %s\n", cfg.ListenAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting new proxy connection: %s\n", err.Error())
				continue
			}
			go func() {
				defer conn.Close()
				if err = socksServer.ServeConn(conn); err != nil {
					log.Printf("error socks proxying a connection: %s\n", err.Error())
				}
			}()
		}
	}()

	return nil
}

func Stop() {
}
