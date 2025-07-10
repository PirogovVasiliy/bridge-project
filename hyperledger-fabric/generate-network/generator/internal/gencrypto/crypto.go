package gencrypto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Org struct {
	Name   string
	Domain string
	Users  int
}

const (
	templatePath = "../templates/crypto-config.tmpl"
	outputDir    = "../../organizations/cryptogen"
)

func GenerateCryptoConfig(orgs []Org) error {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		dst := filepath.Join(
			outputDir,
			fmt.Sprintf("crypto-config-%s.yaml", strings.ToLower(org.Name)))

		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		if err := tpl.Execute(f, org); err != nil {
			return err
		}
		f.Close()
		fmt.Printf("📦  %s создан\n", dst)
	}
	return nil
}
