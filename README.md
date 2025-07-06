Requires Postgres and Go.

Install gator CLI with the command 'go install' in the gator/ directory. 

Requires a .gatorconfig.json config file in the home directory with the following information:
{
    "db_url" : "postgres://user:password@localhost:5432/gator", // Path to your postgres database.
    "current_user_name" : "" 
}

All commands start with gator followed by the desired command. Ex: gator feeds. 
Commands:
1. register - Registers a user. Logs the user in automatically after registering. 
2. login - Logs in with specified user. 
3. users - Returns all registered users. 
4. addfeed - Adds rss feed to the database. If feed is already owned by another users, will follow the feed instead. 
5. feeds - Returns all registered rss feeds. 
6. following - Returns a list of all feeds the current user either owns or follows. 
7. unfollow - Unfollows specified rss. If owner unfollows, feed ownership will transfer to the longest follower of the feed if another follower exists. 
8. browse - Returns the contents of the most recent feeds that the current user follows. Number of feeds returned specified as an argument. 
9. agg - Starts aggregating rss feeds at the specified interval. Ctrl+C to interrupt. 