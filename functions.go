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
	cmds.list["addfeed"] = handlerAddFeed
	cmds.list["feeds"] = handlerFeeds
	cmds.list["follow"] = handlerFollow
	cmds.list["following"] = handlerFollowing
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
		return fmt.Errorf("argumment number error: expected 1. receved %d.\nexpected syntax: login [name]", argNum)
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
		return fmt.Errorf("argumment number error: expected 1. receved %d.\nexpected syntax: register [name]", argNum)
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
	argNum := len(cmd.arguments)
	if argNum > 0 {
		return fmt.Errorf("error: reset command takes no arguments")
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Users have been reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum > 0 {
		return fmt.Errorf("error: users command takes no arguments")
	}

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

func handlerAddFeed(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum != 2 {
		return fmt.Errorf("argumment number error: expected 2. receved %d.\nexpected syntax: addfeed \"[name]\" [url]", argNum)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feedName := cmd.arguments[0]
	feedURL := cmd.arguments[1]

	_, err = fetchFeed(context.Background(), feedURL) // Makes sure url is good
	if err != nil {
		return err
	}

	params := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	}

	feed, err := s.db.AddFeed(context.Background(), params)
	if err != nil {
		if fmt.Sprint(err) == "pq: duplicate key value violates unique constraint \"feeds_url_key\"" {
			fmt.Println("Another user already owns this feed!")
			fmt.Println("Attempting to follow feed...")

			var newArguments []string
			newArguments = append(newArguments, cmd.arguments[1])
			cmd.arguments = newArguments
			err = handlerFollow(s, cmd)
			if err != nil {
				return err
			}
			return nil

		} else {
			return err
		}
	}

	fmt.Println("feed successfully added")
	fmt.Println("ID:", feed.ID)
	fmt.Println("CreatedAt:", feed.CreatedAt)
	fmt.Println("UpdatedAt:", feed.UpdatedAt)
	fmt.Println("Name:", feed.Name)
	fmt.Println("url:", feed.Url)
	fmt.Println("UserID:", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum > 0 {
		return fmt.Errorf("error: feeds command takes no arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("no feeds found")
		return nil
	}

	for i, feed := range feeds {
		fmt.Printf("%d:\n", i+1)
		fmt.Println("Feed Name:", feed.Name)
		fmt.Println("URL:", feed.Url)
		fmt.Println("Owner Name:", feed.Name_2)
		fmt.Println("")
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum != 1 {
		return fmt.Errorf("argumment number error: expected 1. receved %d.\nexpected syntax: follow [url]", argNum)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feedURL := cmd.arguments[0]

	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Follow Successful")
	fmt.Println(currentUser.Name, "now following", feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	argNum := len(cmd.arguments)
	if argNum > 0 {
		return fmt.Errorf("error: following command takes no arguments")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feedsFollowing, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	feedsOwned, err := s.db.GetFeedsOwned(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	if len(feedsFollowing) == 0 && len(feedsOwned) == 0 {
		fmt.Println("you do not own or follow any feeds")
	}

	fmt.Printf("Feeds followed by %s:\n", currentUser.Name)
	ct := 1
	for _, row := range feedsFollowing {
		fmt.Printf("%d. \"%s\" owned by %s\n", ct, row.Name, row.Name_2)
		ct++
	}
	for _, row := range feedsOwned {
		fmt.Printf("%d. \"%s\" owned by you\n", ct, row)
	}

	return nil
}
