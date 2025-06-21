package main

import (
	"fmt"
	"os"

	"github.com/PharmacyDoc2018/gator/internal/config"
)

func printConfigFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	filePath := homeDir + "/.gatorconfig.json"
	configFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	contents := make([]byte, 64)
	configFile.Read(contents)
	configFile.Close()
	fmt.Println(string(contents))

}

func main() {
	printConfigFile()

	fmt.Printf("\n")
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	configStruct.SetUser("Jessica")
	fmt.Println(configStruct)
	printConfigFile()

	fmt.Printf("\n")
	configStruct.SetUser("Joseph")
	fmt.Println(configStruct)
	printConfigFile()

}
