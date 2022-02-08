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
	}
)

func arrContains[Type comparable](search Type, array []Type) *Type {
	for _, element := range array {
		if element == search {
			return &element
		}
	}
	return nil
}

func mapContains[Key comparable, Value comparable]
	(search Value, source map[Key]Value) *Key {
	for key, value := range source {
		if value == search {
			return &key
		}
	}
	return nil
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

func resp(content string, ephemeral bool) dgo.InteractionResponse {
	flags := 0
	if ephemeral {
		flags := 1 << 6
	}

	return dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
			Flags:   flags,
		},
	}
}
