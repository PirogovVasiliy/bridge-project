package gencompose

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"generator/internal/gencrypto"
)

type orgCompose struct {
	Name, Domain, MSP string
	PeerPort          int
	ChaincodePort     int
	OpsPort           int
}

const (
	tplDir    = "../templates"
	outputDir = "../../compose"
)

func enrich(orgs []gencrypto.Org) []orgCompose {
	composeOrgs := make([]orgCompose, len(orgs))
	for i, o := range orgs {
		idx := i + 1
		p := 7051 + (idx-1)*2000 // 7051, 9051 …
		composeOrgs[i] = orgCompose{
			Name:          o.Name,
			Domain:        o.Domain,
			MSP:           fmt.Sprintf("%sMSP", o.Name),
			PeerPort:      p,
			ChaincodePort: p + 1,
			OpsPort:       9444 + (idx - 1), // 9444, 9445…
		}
	}
	return composeOrgs
}

func Generate(orgs []gencrypto.Org) error {
	orcs := enrich(orgs)

	if err := render(
		filepath.Join(tplDir, "compose-test-net.tmpl"),
		filepath.Join(outputDir, "compose-test-net.yaml"),
		orcs); err != nil {
		return err
	}

	if err := render(
		filepath.Join(tplDir, "docker-compose-test-net.tmpl"),
		filepath.Join(outputDir, "docker", "docker-compose-test-net.yaml"),
		orcs); err != nil {
		return err
	}

	return nil
}

func render(tplPath, dstPath string, orcs []orgCompose) error {

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct{ Orgs []orgCompose }{orcs}
	if err := tpl.Execute(f, data); err != nil {
		return err
	}

	fmt.Printf("📦  %s создан\n", dstPath)
	return nil
}
