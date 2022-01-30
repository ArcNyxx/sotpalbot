package main

import (
	"math/rand"

	dgo "github.com/bwmarrin/discordgo"
)

func start(ss *dgo.Session, in *dgo.InteractionCreate) {
	trustedRole, untrustedRole := "", ""
	for _, role := range ss.GuildRoles(guild) {
		if role.Name == "SOTPAL Trusted" {
			trustedRole = role.ID
			if untrustedRole != "" {
				break
			}
		} else if role.Name == "SOTPAL Untrusted" {
			untrustedRole = role.ID
			if trustedRole != "" {
				break
			}
		}
	}
	if trustedRole == "" || untrustedRole == "" {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "Both the \"SOTPAL Trusted\" and " +
					"\"SOTPAL Untrusted\" roles do not exist",
				Flags:   1 << 6, // Ephemeral
			},
		})
		return
	}

	if _, ok := state[in.GuildID]; ok {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "A game of SOTPAL is already running",
				Flags:   1 << 6, // Ephemeral
			},
		})
		return
	}

	if isTrusted(&in.Member, false) {
		state[in.GuildID] = State{
			TrustedRole:    trustedRole,
			UntrustedRole:  untrustedRole,
			TrustedUsers:   make([]string),
			UntrustedUsers: make([]string),
			Submissions:    make(map[string]string),
		}

		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "<@" + in.Member.User.id + "> has " +
					"started a new game of SOTPAL",
			},
		})
	} else {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You do not have the \"SOTPAL Trusted\" role",
				Flags:   1 << 6, // Ephemeral
			},
		})
	}
}

func end(ss *dgo.Session, in *dgo.InteractionCreate) {
	if gameInactive(ss, in) {
		return
	}

	if isTrusted(&in.Member, false) {
		delete(state, ss.GuildID)
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "<@" + in.Member.User.id + "> has " +
					"ended the current game of SOTPAL",
			},
		})
	} else {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You do not have the \"SOTPAL Trusted\" role",
				Flags:   1 << 6, // Ephemeral
			},
		})
	}
}

func article(ss *dgo.Session, in *dgo.InteractionCreate) {
	if gameInactive(ss, in) {
		return
	}

	if isTrusted(&in.Member, true) {
		
	} else {
		ss.In
	}
}
