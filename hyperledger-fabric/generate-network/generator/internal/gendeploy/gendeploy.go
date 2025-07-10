package gendeploy

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"generator/internal/gencompose"
)

const (
	tplPath = "../templates/deploy-token.tmpl"
	outPath = "../../deployToken.sh"
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

	for i := range orgs {
		orgs[i].Name = strings.ToLower(orgs[i].Name)
	}

	if err := tpl.Execute(f, struct{ Orgs []gencompose.OrgCompose }{orgs}); err != nil {
		return err
	}

	if err := os.Chmod(outPath, 0o755); err != nil {
		return err
	}
	fmt.Println("📦  deployToken.sh создан")
	return nil
}
