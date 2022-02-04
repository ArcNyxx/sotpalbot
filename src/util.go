package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func isTrusted(member *dgo.Member) bool {
	for _, role := range member.Roles {
		if role == state[member.GuildID].trustedRole {
			return true
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
	return false
}

func isPlayer(ID string, guildID string) bool {
	for key, _ := range state[guildID].Submissions {
		if key == ID {
			return true
		}
	}
	return false
}
