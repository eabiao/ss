package main

import (
	"log"
	"net"
)

func main() {
	ssk, err := net.Listen("tcp", config.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen:", config.Listen)

	for {
		sk, err := ssk.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnect(sk)
	}
}
