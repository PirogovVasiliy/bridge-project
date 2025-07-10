package envvar

import (
	"fmt"
	"generator/internal/gencompose"
	"os"
	"strconv"
)

const path = "../../scripts/internal/envVar.env"

func GenerateEnvVar(orgs []gencompose.OrgCompose) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "ORG_COUNT=%d\n", len(orgs))

	var domains, msp, peerPorts, chaincodePorts, opsPorts, ca []string
	for _, o := range orgs {
		domains = append(domains, o.Domain)
		msp = append(msp, o.MSP)
		peerPorts = append(peerPorts, strconv.Itoa(o.PeerPort))
		chaincodePorts = append(chaincodePorts, strconv.Itoa(o.ChaincodePort))
		opsPorts = append(opsPorts, strconv.Itoa(o.OpsPort))
		ca = append(ca, fmt.Sprintf("organizations/peerOrganizations/%s/tlsca/tlsca.%s-cert.pem", o.Domain, o.Domain))
	}

	printArr(f, "ORG_DOMAINS", domains)
	printArr(f, "MSP_IDS", msp)
	printArr(f, "PEER_PORTS", peerPorts)
	printArr(f, "CHAINCODE_PORTS", chaincodePorts)
	printArr(f, "OPS_PORTS", opsPorts)
	printArr(f, "TLS_CA_PATHS", ca)
	fmt.Printf("📦  %s создан\n", path)
	return nil
}

func printArr(f *os.File, name string, atributes []string) {
	fmt.Fprintf(f, "%s=(", name)
	for _, v := range atributes {
		fmt.Fprintf(f, "\"%v\" ", v)
	}
	fmt.Fprintln(f, ")")
}
