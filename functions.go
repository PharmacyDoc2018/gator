package main

import (
	"database/sql"
	"fmt"

	"github.com/PharmacyDoc2018/gator/internal/config"
	"github.com/PharmacyDoc2018/gator/internal/database"
)

func initNewState() (*state, error) {
	var newState state
	config, err := config.Read()
	if err != nil {
		return nil, err
	}
	newState.config = config

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator")
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)
	newState.db = dbQueries

	return &newState, nil
}

func initCommands() *commands {
	var cmds commands
	cmds.list = make(map[string]func(*state, command) error)
	cmds.list["login"] = handlerLogin
	return &cmds
}

func exeCommand(s *state, c *commands, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments")
	}

	var cmd command
	cmd.name = args[1]
	cmd.arguments = args[2:]

	err := c.run(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum != 1 {
		return fmt.Errorf("argumment number error. expected 1 argument. receved %d", argNum)
	}

	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil

}
