package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func submitcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if arrContains(state[in.GuildID].TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &UntrustedUser)
		return
	}

	article := in.ApplicationCommandData().Options[0].StringValue()
	content := "You have set your article to \"" + article + "\"."
	if _, ok := state[in.GuildID].Submissions[in.Member.User.ID]; ok {
		content = content[:len(content) - 1] + ", overwriting \"" +
			state[in.GuildID].Submissions[in.Member.User.ID] + "\"."
	}
	state[in.GuildID].Submissions[in.Member.User.ID] = article
	ss.InteractionRespond(in.Interaction, err(content))
}

func removecmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if len(in.ApplicationCommandData().Options) == 0 {
		if arrContains(state[in.GuildID].UntrustedRole, in.Member.Roles) == nil {
			ss.InteractionRespond(in.Interaction, &UntrustedUser)
			return
		}

		if _, ok := state[in.GuildID].Submissions[in.Member.User.ID]; !ok {
			ss.InteractionRespond(in.Interaction, err(
				"You have not submitted an article."))
			return
		}

		ss.InteractionRespond(in.Interaction, err(
			"You have removed your own article \"" +
			state[in.GuildID].Submissions[in.Member.User.ID] +
			"\"."))
		delete(state[in.GuildID].Submissions, in.Member.User.ID)
		return
	}

	if arrContains(state[in.GuildID].TrustedRole, in.Member.Roles) == nil {
		ss.InteractionRespond(in.Interaction, &NonTrustedUser)
		return
	}

	article := in.ApplicationCommandData().Options[0].StringValue()
	untrust := in.ApplicationCommandData().Options[1].BoolValue()

	if player := mapContains(article, state[in.GuildID].Submissions);
		player != nil {
		content := "<@" + in.Member.User.ID + "> has removed the " +
			"article \"" + article + "\"."
		if untrust {
			content = content[:len(content) - 1] + " and " +
				"untrusted <@" + *player + ">"
			if ofs := ss.GuildMemberRoleAdd(in.GuildID, in.Member.User.ID,
				state[in.GuildID].UntrustedRole); ofs != nil {
				ss.InteractionRespond(in.Interaction, err(
					"Unable to give <@" + *player + "> the " +
						"\"SOTPAL Untrusted\" role."))
			}
		}
		delete(state[in.GuildID].Submissions, *player)
		ss.InteractionRespond(in.Interaction, resp(content))
	} else {
		ss.InteractionRespond(in.Interaction, err(
			"\"" + article + "\" is not the name of a submitted " +
			"article."))
	}
}

func printcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}

	if len(state[in.GuildID].Submissions) == 0 {
		ss.InteractionRespond(in.Interaction, &NoSubmissions)
		return
	}

	ss.InteractionRespond(in.Interaction, resp(
		"<@" + in.Member.User.ID + "> has requested the list of " +
		"submitted articles: " +
		mentionSubmit(state[in.GuildID].Submissions, false)))
}

func clearcmd(ss *dgo.Session, in *dgo.InteractionCreate) {
	if _, ok := state[in.GuildID]; !ok {
		ss.InteractionRespond(in.Interaction, &GameNotRunning)
		return
	}
	
	if len(state[in.GuildID].Submissions) == 0 {
		ss.InteractionRespond(in.Interaction, &NoSubmissions)
		return
	}

	newstate := state[in.GuildID]
	newstate.Submissions = make(map[string]string)
	state[in.GuildID] = newstate
	ss.InteractionRespond(in.Interaction, resp(
		"<@" + in.Member.User.ID + "> has cleared the list of " +
			"submitted articles."))
}
