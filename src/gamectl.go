package main

import (
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func startcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; ok {
		ss.InteractionRespond(in.Interaction, err(
			"A game of SOTPAL is already active."))
		return
	}

	if roles, err = ss.GuildRoles(in.GuildID); err != nil {
		ss.InteractionRespond(in.Interaction, err(
			"Unable to get server roles."))
		return
	}
	
	var state State
	for _, role := range roles {
		if role.Name == "SOTPAL Trusted" {
			state.TrustedRole = role.ID
			if state.UntrustedRole != "" {
				break
			}
		} else if role.Name == "SOTPAL Untrusted" {
			state.UntrustedRole = role.ID
			if state.TrustedRole != "" {
				break
			}
		}
	}
	if state.TrustedRole == "" || state.UntrustedRole == "" {
		ss.InteractionRespond(in.Interaction, err(
			"The required \"SOTPAL Trusted\" and \"SOTPAL " +
			"Untrusted\" roles do not exist."))
		return
	}

	if arrContains(state.TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
	} else {
		state.Submissions = make(map[string]string)
		ss.InteractionRespond(in.Interaction, resp(
			"<@" + in.Member.User.ID + "> has started a new " +
			"game of SOTPAL!"))
	}
}

func endcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if arrContains(state.TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
	} else {
		delete(state, ss.GuildID)
		ss.InteractionRespond(in.Interaction, resp(
			"<@" + in.Member.User.ID + "> has ended the game of " +
			"SOTPAL!"
		))
	}
}

func articlecmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if arrContains(state.TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host != "" {
		ss.InteractionRespond(in.Interaction, err(
			"A round of SOTPAL is already active."))
		return
	}

	if len(state[in.GuildID].Submissions) < 2 {
		ss.InteractionRespond(in.Interaction, err(
			"Fewer than two players are currently available."))
		return
	}

	if _, ok := state[in.GuildID].Submissions[in.Member.User.ID]; !ok {
		ss.InteractionRespond(in.Interaction, err(
			"You cannot submit an article and host a game. " +
			"Please remove your article and try again."))
		return
	}

	cnt, random := 0, rand.Intn(len(state[in.GuildID].Submissions))
	for player, article := range state[in.GuildID].Submissions {
		if cnt == random {
			state[in.GuildID].Player, state[in.GuildID].Article =
				player, article
			delete(state[in.GuildID].Submissions, player)
			break
		}
		cnt++
	}
	state[in.GuildID].Host = in.Member.User.ID

	ss.InteractionRespond(in.Interaction, resp(
		"<@" + in.Member.User.ID + "> has started a new round of " +
		"SOTPAL! The article is \"" + state[in.GuildID].Article +
		"\", and the players are " + 
		mentionSubmit(state[in.GuildID].Submissions, true) + "."))
}

func guesscmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if arrContains(state.TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host == "" {
		ss.InteractionRespond(in.Interaction, err(
			"A round of SOTPAL is not currently active."))
		return
	}
	
	if state[in.GuildID].Host != in.Member.User.ID {
		ss.InteractionRespond(in.Interaction, err(
			"You are not the host of the current round of SOTPAL."))
		return
	}

	guess := in.ApplicationCommandData().Options[0].UserValue(nil).ID
	content := "<@" + state[in.GuildId].Host + "> guessed that <@" + guess +
		"> submitted the article \"" + state[in.GuildID].Article + "\""
	if guess == state[in.GuildID].Player {
		content += " and was correct!"
	} else {
		if _, ok := state[in.GuildID].Submissions; !ok {
			ss.InteractionRespond(in.Interaction, err(
				"<@" + guess + "> is not playing in this round."))
			return
		}
		content += ", but it was actually <@" + state[in.GuildID].Player +
			">!"
	}
	state[in.GuildID].Host, state[in.GuildID].Player,
		state[in.GuildID].Article = "", "", ""
	ss.InteractionRespond(in.Interaction, resp(content))
}
