package main

import (
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func start(ss *dgo.Session, in *dgo.InteractionCreate) {
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
						"started a new game of SOTPAL.",
				},
			})
			return
		}
	}
	ss.InteractionResponse(in.Interaction, &NonTrustedUser)
}

func end(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member, false) {
		ss.InteractionResponse(in.Interaction, &NonTrustedUser)
		return
	}

	delete(state, ss.GuildID)
	ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "<@" + in.Member.User.id + "> has " +
				"ended the current game of SOTPAL.",
		},
	})
}

func article(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member, true) {
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
	}
	
	if len(state[in.GuildID].Submissions) < 2 {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "Fewer than two players are currently available.",
				Flags:   1 << 6,
			},
		})
	}

	for key, _ := range state.Submissions {
		if key == in.Member.User.ID {
			ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
				Type: InteractionResponseChannelMessageWithSource,
				Data: &dgo.InteractionResponseData{
					Content: "You have submitted an article and attempted " +
						"to start a round. Please remove your article " +
						"and try again.",
					Flags:   1 << 6,
				},
			})
		}
	}

	mentionString := ""
	state[in.GuildID].Host = in.Member.User.ID

	random := rand.Intn(len(state.Submissions))
	for key, value := range state.Submissions {
		if random == 0 {
			state[in.GuildID].Player  = key
			state[in.GuildID].Article = value
			delete(state, key)
		}
		
		mentionString += "<@" + key + ">, "
		random -= 1
	}

	ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "New round started by <@" + in.Member.User.ID + ">! " +
				"The article is " + state[in.GuildID].Article + ", and the " +
				"players are " + mentionString[:len(mentionString) - 2] + ".",
		}
	})
}

func guess(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionResponse(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member, true) {
		ss.InteractionResponse(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host != in.Member.User.ID {
		ss.InteractionResponse(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "You are not the host of a currently "
					"active round of SOTPAL.",
				Flags:   1 << 6,
			},
		})
	}

	user := in.ApplicationCommandData().Options[0].UserValue(nil).ID
	if user == state[in.GuildID].Player {
		content := "<@" + state[in.GuildId].Host + "> guessed that <@" +
			state[in.GuildID].Player + "> submitted the article " +
			state[in.GuildID].Article + ", and was correct!"
	} else {
		content := "<@" + state[in.GuildID].Host + "> guessed that <@" +
			user + "> submitted the article " + state[in.GuildID].Article +
			", but it was actually <@" + state[in.GuildID].Article + ">!"
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
