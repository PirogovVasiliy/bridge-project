package genconfigtx

import (
	"fmt"
	"os"
	"text/template"

	"generator/internal/gencompose"
)

const (
	tplPath = "../templates/configtx.tmpl"
	outPath = "../../configtx/configtx.yaml"
)

func Generate(orgs []gencompose.OrgCompose) error {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tpl.Execute(f, struct{ Orgs []gencompose.OrgCompose }{orgs}); err != nil {
		return err
	}

	fmt.Printf("📦  %s создан\n", outPath)
	return nil
}
