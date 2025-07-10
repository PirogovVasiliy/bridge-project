package gencompose

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"generator/internal/gencrypto"
)

type OrgCompose struct {
	Name, Domain, MSP string
	PeerPort          int
	ChaincodePort     int
	OpsPort           int
}

const (
	tplDir    = "../templates"
	outputDir = "../../compose"
)

func enrich(orgs []gencrypto.Org) []OrgCompose {
	composeOrgs := make([]OrgCompose, len(orgs))
	for i, o := range orgs {
		idx := i + 1
		p := 7051 + (idx-1)*2000 // 7051, 9051 …
		composeOrgs[i] = OrgCompose{
			Name:          o.Name,
			Domain:        o.Domain,
			MSP:           fmt.Sprintf("%sMSP", o.Name),
			PeerPort:      p,
			ChaincodePort: p + 1,
			OpsPort:       9444 + (idx - 1), // 9444, 9445…
		}
	}
	//composeOrgs[0].MSP = "Org1MSP"
	//composeOrgs[1].MSP = "Org2MSP"
	return composeOrgs
}

func Generate(orgs []gencrypto.Org) ([]OrgCompose, error) {
	composeOrgs := enrich(orgs)

	if err := render(
		filepath.Join(tplDir, "compose-test-net.tmpl"),
		filepath.Join(outputDir, "compose-test-net.yaml"),
		composeOrgs); err != nil {
		return nil, err
	}

	if err := render(
		filepath.Join(tplDir, "docker-compose-test-net.tmpl"),
		filepath.Join(outputDir, "docker", "docker-compose-test-net.yaml"),
		composeOrgs); err != nil {
		return nil, err
	}

	return composeOrgs, nil
}

func render(tplPath, dstPath string, orcs []OrgCompose) error {

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct{ Orgs []OrgCompose }{orcs}
	if err := tpl.Execute(f, data); err != nil {
		return err
	}

	fmt.Printf("📦  %s создан\n", dstPath)
	return nil
}
