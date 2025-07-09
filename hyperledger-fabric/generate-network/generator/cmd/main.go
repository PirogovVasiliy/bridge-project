package main

import (
	"fmt"
	"generator/internal/mycrypto"
	"log"
)

func main() {
	orgs := []mycrypto.Org{
		{Name: "logitech", Domain: "Logitech.com", Users: 3},
		{Name: "zetgaming", Domain: "zetgaming.com", Users: 1},
	}

	fmt.Println(orgs[0])
	fmt.Println(orgs[1])

	if err := mycrypto.GenerateCryptoConfig(orgs); err != nil {
		log.Fatalf("ошибка генерации: %v", err)
	}
}
