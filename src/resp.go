package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
	GameNotRunning = dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "A game of SOTPAL is not currently active."
			Flag:    1 << 6,
		},
	}

	NonTrustedUser = dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "You are not a trusted user (lacking \"SOTPAL " +
				"Trusted\" role).",
			Flags:   1 << 6,
		},
	}

	UntrustedUser = dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "You are an untrusted user (having the \"SOTPAL " +
				"Untrusted\" role).",
			Flags:   1 << 6,
		},
)
