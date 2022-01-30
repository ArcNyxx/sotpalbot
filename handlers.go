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

	TrustedUsers   []string
	UntrustedUsers []string
}

var (
	state = make(map[string]State)

	handlers = map[string]func(ss *dgo.Session, in *dgo.InteractionCreate){
		"start": start,
		"end":   end,
	}
)
