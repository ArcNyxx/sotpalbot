package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func gameInactive(ss *dgo.Session, in *dgo.InteractionCreate) bool {
	if _, ok := state[in.GuildID]; ok {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "A of SOTPAL is not currently running",
				Flags:   1 << 6, // Ephemeral
			},
		})
		return true
	}
	return false
}

func isTrusted(member *dgo.Member, useState bool) bool {
	for _, role := range member.Roles {
		if role == state[member.GuildID].trustedRole {
			return true
		}
	}
	if useState {
		for _, user := range state[member.GuildID].trustedUsers {
			if user == member.User.ID {
				return true
			}
		}
	}
	return false
}

func isUntrusted(member *dgo.Member) bool {
	for _, role := range member.Roles {
		if role == state[member.GuildID].untrustedRole {
			return true
		}
	}
	for _, user := range state[member.GuildID].untrustedUsers {
		if user == member.User.ID {
			return true
		}
	}
	return false
}
