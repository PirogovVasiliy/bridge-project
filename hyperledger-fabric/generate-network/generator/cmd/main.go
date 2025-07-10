package main

import (
	"generator/internal/envvar"
	"generator/internal/gencompose"
	"generator/internal/genconfigtx"
	"generator/internal/gencrypto"
	"log"
)

func main() {

	orgs := []gencrypto.Org{
		{Name: "logitech", Domain: "logitech.com", Users: 3},
		{Name: "zetgaming", Domain: "zetgaming.com", Users: 1},
		{Name: "acer", Domain: "acer.com", Users: 1},
		{Name: "axnot", Domain: "axnot.com", Users: 2},
	}

	if err := gencrypto.GenerateCryptoConfig(orgs); err != nil {
		log.Fatalln("generate crypto-config error", err)
	}

	composeOrgs, err := gencompose.Generate(orgs)

	if err != nil {
		log.Fatalln("generate compose error", err)
	}

	err = envvar.GenerateEnvVar(composeOrgs)
	if err != nil {
		log.Fatalln("generate envVar error", err)
	}

	err = genconfigtx.Generate(composeOrgs)
	if err != nil {
		log.Fatalln("generate configtx error", err)
	}
}
