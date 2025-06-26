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

func initNewState() (*state, *sql.DB, error) {
	var newState state
	config, err := config.Read()
	if err != nil {
		return nil, nil, err
	}
	newState.config = config

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator")
	if err != nil {
		return nil, nil, err
	}
	dbQueries := database.New(db)
	newState.db = dbQueries

	return &newState, db, nil
}

func initCommands() *commands {
	var cmds commands
	cmds.list = make(map[string]func(*state, command) error)
	cmds.list["login"] = handlerLogin
	cmds.list["register"] = handlerRegister
	cmds.list["reset"] = handlerReset
	cmds.list["users"] = handlerUsers
	cmds.list["agg"] = handlerAgg
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
		if fmt.Sprint(err) == "pq: duplicate key value violates unique constraint \"users_name_key\"" {
			return fmt.Errorf("error: user already exists")
		} else {
			return err
		}
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

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Users have been reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("no users registered")
		return nil
	}
	var printLine string
	for _, user := range users {
		printLine = fmt.Sprintf("* %s", user)
		if user == s.config.CurrentUserName {
			printLine += " (current)"
		}
		fmt.Println(printLine)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(rss)
	return nil
}
