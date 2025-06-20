package main

import (
	"fmt"

	"github.com/PharmacyDoc2018/gator/internal/config"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	gatorConfig.SetUser("Joseph")

	fmt.Println(gatorConfig.DbURL)
	fmt.Println(gatorConfig.CurrentUserName)
}
