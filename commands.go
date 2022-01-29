package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

var commands = []*dgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Start the game of SOTPAL, " +
			"creating a new game state for the guild",
	},
	{
		Name:        "end",
		Description: "End the game of SOTPAL, " +
			"removing any saved game state for the guild",
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
				Name:        "player",
				Description: "The player who submitted the article",
				Required:    true
			},
		},
	},

	{
		Name:        "trust",
		Description: "Toggle host permissions for a player, " +
			"temporarily or permanently",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionUser,
				Name:        "player",
				Description: "The player whose permissions to toggle",
				Required:    true
			},
			{
				Type:        dgo.ApplicationCommandOptionBool,
				Name:        "perm",
				Description: "Whether to toggle permissions in the context " +
					"of the game state or permanently",
				Required:    false
			},
		},
	},
	{
		Name:        "untrust",
		Description: "Toggle submission permissions for a player, " +
			"temporarily or permanently",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionUser,
				Name:        "player",
				Description: "The player whose permissions to toggle",
				Required:    true,
			},
			{
				Type:        dgo.ApplicationCommandOptionBool,
				Name:        "perm",
				Description: "Whether to toggle permissions in the context " +
					"of the game or permanently",
				Required:    false,
			},
		},
	},

	{
		Name:        "submit",
		Description: "Submit an article to the list",
		Options:     []*dgo.ApplicationCommandOption{
			{
				Type:        dgo.ApplicationCommandOptionString,
				Name:        "article",
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
				Name:        "article",
				Description: "The article to remove from the list",
				Required:    false,
			},
			{
				Type:        dgo.ApplicationCommandOptionBool,
				Name:        "print",
				Description: "Whether to print the player who submitted the article",
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
