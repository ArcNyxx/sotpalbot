package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func submitcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	article := in.ApplicationCommandData().Options[0].StringValue()
	if _, ok := state[in.GuildID].Submissions[in.Member.User.ID]; ok {
		content := "You have set your article to " + article + ", overwriting " +
		state[in.GuildID].Submissions[in.Member.User.ID] + "."
	} else {
		content := "You have set your article to " + article + "."
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
	
}
