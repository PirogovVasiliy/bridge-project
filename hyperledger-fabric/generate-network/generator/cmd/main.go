package main

import (
	"generator/internal/gencompose"
	"generator/internal/gencrypto"
	"log"
)

func main() {

	orgs := []gencrypto.Org{
		{Name: "logitech", Domain: "logitech.com", Users: 3},
		{Name: "zetgaming", Domain: "zetgaming.com", Users: 1},
		{Name: "acer", Domain: "acer.com", Users: 1},
	}

	if err := gencrypto.GenerateCryptoConfig(orgs); err != nil {
		log.Fatalln("generate crypto-config error", err)
	}

	if err := gencompose.Generate(orgs); err != nil {
		log.Fatalln("generate compose error", err)
	}
}
