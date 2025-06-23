package main

import (
	"fmt"
	"os"
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
	state, err := initNewState()
	if err != nil {
		fmt.Println(err)
	}

	commands := initCommands()

	err = exeCommand(state, commands, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
