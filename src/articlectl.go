package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func submitcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if isUntrusted(in.Member) {
		ss.InteractionRespond(in.Interaction, &UntrustedUser)
		return
	}

	article := in.ApplicationCommandData().Options[0].StringValue()
	if _, ok := state[in.GuildID].Submissions[in.Member.User.ID]; ok {
		content := "You have set your article to \"" + article + "\", overwriting \"" +
		state[in.GuildID].Submissions[in.Member.User.ID] + "\"."
	} else {
		content := "You have set your article to \"" + article + "\"."
	}
	state[in.GuildID].Submissions[in.Member.User.ID] = article

	ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	})
}

func removecmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if len(i.ApplicationCommandData().Options) == 1 {
		if !isTrusted(&in.Member) {
			ss.InteractionResponse(in.Interaction, &NonTrustedUser)
			return
		}
		user := in.ApplicationCommandData().Options[0].UserValue(nil).ID
	} else {
		user := in.Member.User.ID
	}

	
}

func printcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if len(state[in.GuildID].Submissions) == 0 {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "No articles have been submitted.",
				Flags:   1 << 6,
			},
		})
		return
	}

	i, length := 0, len(state[in.GuildID].Submissions
	articles := ""
	for _, value := range state[in.GuildID].Submissions {
		articles += "\"" + value + "\""
		if i != length - 1 {
			articles += ", "
			if i == length - 2 {
				articles += "and "
			}
		}
	}

	ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseCHannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "<@" + in.Member.User.ID + "> has requested the list " +
				"of submitted articles: " + articles,
		},
	})
}

func clearcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}
	
	if len(state[in.GuildID].Submissions) == 0 {
		ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "No articles have been submitted.",
				Flags:   1 << 6,
			},
		})
		return
	}

	state[in.GuildID].Submissions = make(map[string]string)
	ss.InteractionRespond(in.Interaction, &dgo.InteractionResponse{
		Type: InteractionResponseCHannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: "<@" + in.Member.User.ID + "> has cleared "
				"the list of submitted articles",
		},
	})
}
