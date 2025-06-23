package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/PharmacyDoc2018/gator/internal/config"
	"github.com/PharmacyDoc2018/gator/internal/database"
	"github.com/google/uuid"
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
	cmds.list["register"] = handlerRegister
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

	loginName := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), loginName)
	if err != nil {
		return fmt.Errorf("error: user not found")
	}

	err = s.config.SetUser(loginName)
	if err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil

}

func handlerRegister(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum != 1 {
		return fmt.Errorf("argumment number error. expected 1 argument. receved %d", argNum)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	newUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.config.SetUser(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Println(newUser.Name, "was created")
	fmt.Println("ID:", newUser.ID)
	fmt.Println("CreatedAt:", newUser.CreatedAt)
	fmt.Println("UpdatedAt:", newUser.UpdatedAt)
	fmt.Println("Name:", newUser.Name)
	return nil
}
