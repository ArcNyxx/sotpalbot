package main

import (
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func startcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; ok {
		ss.InteractionRespond(in.Interaction, resp(
			"A game of SOTPAL is already active", true
		))
		return
	}

	// Find trusted and untrusted roles; Return if not found
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
		ss.InteractionRespond(in.Interaction, resp(
			"The required \"SOTPAL Trusted\" and \"SOTPAL " +
			"Untrusted\" roles do not exist.", true
		))
		return
	}

	// Create entry if caller has trusted role
	for _, role := range in.Member.Roles {
		if role == trustedRole {
			state[in.GuildID] = State{
				TrustedRole:   trustedRole,
				UntrustedRole: untrustedRole,
				Submissions:   make(map[string]string),
			}

			ss.InteractionRespond(in.Interaction, resp(
				"<@" + in.Member.User.ID + "> has started " +
				"a new game of SOTPAL!", false
			))
			return
		}
	}
	ss.InteractionRespond(in.Interaction, &NonTrustedUser)
}

func endcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	delete(state, ss.GuildID)
	ss.InteractionRespond(in.Interaction, resp(
		"<@" + in.Member.User.ID + "> has ended the game of SOTPAL!",
		false
	))
}

func articlecmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host != "" {
		ss.InteractionRespond(in.Interaction, resp(
			"A round of SOTPAL is already active.", true
		))
		return
	}
	
	if len(state[in.GuildID].Submissions) < 2 {
		ss.InteractionRespond(in.Interaction, resp(
			"Fewer than two players are currently available.", true
		))
		return
	}

	if isPlayer(in.Member.User.ID, in.GuildID) {
		ss.InteractionRespond(in.Interaction, resp(
			"You cannot submit an article and host a game. " +
			"Please remove your article and try again.", true
		))
		return
	}

	mentions := enumSubmissions(in.GuildID, true)
	i, random := 0, rand.Intn(len(state[in.GuildID].Submissions))
	for player, article := range state[in.GuildID].Submissions {
		if i == random {
			state[in.GuildID].Host = in.Member.User.ID
			state[in.GuildID].Player = player
			state[in.GuildID].Article = article
			delete(state[in.GuildID].Submissions, player)
			break
		}
	}

	ss.InteractionRespond(in.Interaction, resp(
		"<@" + in.Member.User.ID + "> has started a new round of " +
		"SOTPAL! The article is \"" + state[in.GuildID].Article +
		"\", and the players are " + mentions + ".", false
	))
}

func guesscmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if !isTrusted(&in.Member) {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	if state[in.GuildID].Host == "" {
		ss.InteractionRespond(in.Interaction, resp(
			"A round of SOTPAL is not currently active.", true
		))
		return
	}
	
	if state[in.GuildID].Host != in.Member.User.ID {
		ss.InteractionRespond(in.Interaction, resp(
			"You are not the host of the current round of "
			"SOTPAL.", true
		))
		return
	}

	guess = in.ApplicationCommandData().Options[0].UserValue(nil).ID
	if guess == state[in.GuildID].Player {
		content := "<@" + state[in.GuildId].Host + "> guessed that " +
			"<@" + state[in.GuildID].Player + "> submitted the " +
			"article \"" + state[in.GuildID].Article + "\" and " +
			"was correct!"
	} else {
		if !isPlayer(guess, in.GuildID) {
			ss.InteractionRespond(in.Interaction, resp(
				"<@ " + guess + "> is not playing in this " +
				"round.", true
			))
			return
		}
		content := "<@" + state[in.GuildID].Host + "> guessed that <@" +
			user + "> submitted the article \"" + state[in.GuildID].Article +
			"\", but it was actually <@" + state[in.GuildID].Article + ">!"
	}

	state[in.GuildID].Host = ""
	state[in.GuildID].Player = ""
	state[in.GuildID].Article = ""
	ss.InteractionRespond(in.Interaction, resp(content, false))
}
