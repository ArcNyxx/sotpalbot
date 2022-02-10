package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

type State struct {
	Host           string
	Article        string
	Player         string
	
	Submissions    map[string]string

	TrustedRole    string
	UntrustedRole  string
}

var state = make(map[string]State)

var commands = []*dgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Start the game of SOTPAL",
	},
	{
		Name:        "end",
		Description: "End the game of SOTPAL",
	},
	{
		Name:        "article",
		Description: "Start a round of SOTPAL by randomly selecting " +
			"any of the submitted articles",
	},
	{
		Name:        "guess",
		Description: "End a round of SOTPAL by guessing the player" +
			"who submitted the article",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionUser,
				Name:        "Player",
				Description: "The player who is guessed to have " +
					"submitted the article",
				Required:    true,
			},
		},
	},

	{
		Name:        "submit",
		Description: "Submit an article to the list",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionString,
				Name:        "Article",
				Description: "The article to submit to the list",
				Required:    true,
			},
		},
	},
	{
		Name:        "remove",
		Description: "Remove an article from the list",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionString,
				Name:        "Article",
				Description: "The article to remove from the list, " +
					"defaults to own article",
				Required:    false,
			},
			{
				Type:        dgo.ApplicationCommandOptionBoolean,
				Name:        "Untrust",
				Description: "Whether to untrust the submitted of the article",
				Required:    false,
			},
		},
	},
	{
		Name:        "print",
		Description: "Print the list of articles",
	},
	{
		Name:        "clear",
		Description: "Clear the list of articles",
	},
}

var handlers = map[string]func(ss *dgo.Session, in *dgo.InteractionCreate){
	"start":   startcmd,
	"end":     endcmd,
	"article": articlecmd,
	"guess":   guesscmd,

	"submit":  submitcmd,
	"remove":  removecmd,
	"print":   printcmd,
	"clear":   clearcmd,
}
