package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

var (
	GameNotRunning = dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "A game of SOTPAL is not currently active.",
			Flags:   1 << 6,
		},
	}

	NonTrustedUser = dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "You are not a trusted user (lacking \"SOTPAL " +
				"Trusted\" role).",
			Flags:   1 << 6,
		},
	}

	UntrustedUser = dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "You are an untrusted user (having the \"SOTPAL " +
				"Untrusted\" role).",
			Flags:   1 << 6,
		},
	}

	NoSubmissions = dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "No articles have been submitted.",
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

func mapContains[Key comparable, Value comparable](search Value,
	source map[Key]Value) *Key {
	for key, value := range source {
		if value == search {
			return &key
		}
	}
	return nil
}

func mentionSubmit(submissions map[string]string, players bool) string {
	ret, i, length := "", 0, len(submissions)
	for player, article := range submissions {
		if players {
			ret += "<@" + player + ">"
		} else {
			ret += "\"" + article + "\""
		}

		if i != length - 1 && length != 2 {
			ret += ", "
		}
		if i == length - 2 {
			if length == 2 {
				ret += " "
			}
			ret += "and "
		}
		i++
	}
	return ret
}

func resp(content string) *dgo.InteractionResponse {
	return &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
		},
	}
}

func err(content string) *dgo.InteractionResponse {
	return &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	}
}
