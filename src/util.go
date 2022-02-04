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

func isPlayer(playerid string, guild string) bool {
	for player, _ := range state[guild].Submissions {
		if player == playerid {
			return true
		}
	}
	return false
}

func enumSubmissions(guild string, players bool) string {
	ret := ""
	i, length := 0, len(state[guild].Submissions)
	for player, article := range state[guild].Submissions {
		if players {
			ret += "<@" + player + ">"
		} else {
			ret += "\"" + article + "\""
		}

		if i != length - 1 {
			ret += ", "
			if i == length - 2 {
				ret += "and "
			}
		}
		i += 1
	}
	return ret
}
