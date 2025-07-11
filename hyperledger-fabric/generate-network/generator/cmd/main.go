package main

import (
	"generator/internal/envvar"
	"generator/internal/gencompose"
	"generator/internal/genconfigtx"
	"generator/internal/gencrypto"
	"generator/internal/gendeploy"
	"log"
)

func main() {

	orgs := []gencrypto.Org{
		//{Name: "Org1", Domain: "org1.example.com", Users: 3},
		//{Name: "Org2", Domain: "org2.example.com", Users: 1},
		{Name: "Logitech", Domain: "logitech.com", Users: 3},
		{Name: "Zetgaming", Domain: "zetgaming.com", Users: 1},
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

	err = gendeploy.Generate(composeOrgs)
	if err != nil {
		log.Fatalln("generate deployer error", err)
	}
}
