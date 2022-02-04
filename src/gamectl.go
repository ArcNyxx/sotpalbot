package main

import (
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func startcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; ok {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "A game of SOTPAL is currently active.",
				Flags:   1 << 6,
			},
		})
		return
	}

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
				Content: "The required \"SOTPAL Trusted\" and " +
					"\"SOTPAL Untrusted\" roles do not exist.",
				Flags:   1 << 6,
			},
		})
		return
	}

	for _, role := range in.Member.Roles {
		if role == trustedRole {
			state[in.GuildID] = State{
				TrustedRole:   trustedRole,
				UntrustedRole: untrustedRole,
				Submissions:   make(map[string]string),
			}

			ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
				Type: InteractionResponseChannelMessageWithSource,
				Data: &dgo.InteractionResponseData{
					Content: "<@" + in.Member.User.id + "> has " +
						"started a new game of SOTPAL!",
				},
			})
			return
		}
	}
	ss.InteractionResponse(in.Interaction, &NonTrustedUser)
}

func endcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionResponse(in.Interaction, &NonTrustedUser)
		return
	}

	delete(state, ss.GuildID)
	ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "<@" + in.Member.User.id + "> has " +
				"ended the game of SOTPAL!",
		},
	})
}

func articlecmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionResponse(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host != "" {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "A round of SOTPAL is currently active.",
				Flags:   1 << 6,
			},
		})
		return
	}
	
	if len(state[in.GuildID].Submissions) < 2 {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "Fewer than two players are currently available.",
				Flags:   1 << 6,
			},
		})
		return
	}

	for key, _ := range state[in.GuildID].Submissions {
		if key == in.Member.User.ID {
			ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
				Type: InteractionResponseChannelMessageWithSource,
				Data: &dgo.InteractionResponseData{
					Content: "You cannot submit an article and host a game. " +
						"Please remove your article and try again.",
					Flags:   1 << 6,
				},
			})
			return
		}
	}

	i, length := 0, len(state[in.GuildID].Submissions)
	mentions, random := "", rand.Intn(length)
	for key, value := range state[in.GuildID].Submissions {
		if i == random {
			state[in.GuildID].Player = key
			state[in.GuildID].Article = value
			delete(state[in.GuildID].Submissions, key)
		}

		mentions += "<@" + key ">"
		if i != length - 1 {
			mentions += ", "
			if i == length - 2 {
				mentions += "and "
			}
		}
		i += 1
	}
	state[in.GuildID].Host = in.Member.User.ID

	ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "<@" + in.Member.User.ID + "> has started a new round of " +
				"SOTPAL! The article is \"" + state[in.GuildID].Article + "\", " +
				"and the players are " + mentionString + ".",
		}
	})
}

func guesscmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionResponse(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host == "" {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "A round of SOTPAL is not currently active."
				Flags:   1 << 6,
			},
		})
		return
	}
	
	if state[in.GuildID].Host != in.Member.User.ID {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You are not the host of the current round of SOTPAL.",
				Flags:   1 << 6,
			},
		})
		return
	}

	ID = in.ApplicationCommandData().Options[0].UserValue(nil).ID
	if id == state[in.GuildID].Player {
		content := "<@" + state[in.GuildId].Host + "> guessed that <@" +
			state[in.GuildID].Player + "> submitted the article \"" +
			state[in.GuildID].Article + "\" and was correct!"
	} else {
		if !isPlayer(id, in.GuildID) {
			ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
				Type: InteractionResponseChannelMessageWithSource,
				Data: &dgo.InteractionResponseData{
					Content: "<@" + ID + "> is not playing this round",
					Flags:   1 << 6,
				},
			})
		}
		content := "<@" + state[in.GuildID].Host + "> guessed that <@" +
			user + "> submitted the article \"" + state[in.GuildID].Article +
			"\", but it was actually <@" + state[in.GuildID].Article + ">!"
	}

	state[in.GuildID].Host = ""
	state[in.GuildID].Player = ""
	state[in.GuildID].Article = ""

	ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
		},
	})
}
