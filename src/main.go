package main

import (
	"log"
	"os"
	"os/signal"
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func main() {
	if len(os.Args[1:]) != 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		log.Fatalln("usage: sotpal [token]")
	}

	session, err := dgo.New("Bot " + os.Args[1])
	if err != nil {
		log.Fatalln("sotpal: invalid bot token")
	}

	session.AddHandler(func(ss *dgo.Session, rd *dgo.Ready) {
		log.Println("sotpal: bot session created")
		for _, guild := range rd.Guilds {
			if _, err := ss.ApplicationCommandBulkOverwrite(
				ss.State.User.ID, guild.ID, commands); err != nil {
				log.Fatalln("sotpal: unable to register commands")
			}
		}
		log.Println("sotpal: commands registered")
	})
	session.AddHandler(func(ss *dgo.Session, in *dgo.InteractionCreate) {
		if handler, ok := handlers[in.ApplicationCommandData().Name]; ok {
			handler(ss, in)
		}
	})

	if err := session.Open(); err != nil {
		log.Fatalln("sotpal: unable to create bot session")
	}
	defer session.Close()

	rand.Seed(time.Now().UnixNano())

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("sotpal: exiting")
}
