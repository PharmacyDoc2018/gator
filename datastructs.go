package main

import (
	"fmt"

	"github.com/PharmacyDoc2018/gator/internal/config"
	"github.com/PharmacyDoc2018/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	for key, val := range c.list {
		if cmd.name == key {
			err := val(s, cmd)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("command not found")
}

func (c *commands) register(name string, f func(*state, command) error) error {
	for key, _ := range c.list {
		if name == key {
			return fmt.Errorf("command already exists")
		}
	}
	c.list[name] = f
	return nil
}
