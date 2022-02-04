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

var (
	state = make(map[string]State)

	handlers = map[string]func(ss *dgo.Session, in *dgo.InteractionCreate){
		"start":   startcmd,
		"end":     endcmd,
		"article": articlecmd,
		"guess":   guesscmd,

		"submit":  submitcmd,
		"remove":  removecmd,
		"print":   printcmd,
		"clear":   clearcmd,
	}
)
