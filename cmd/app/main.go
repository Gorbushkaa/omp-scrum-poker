package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	handlers "omppoker/internal"
	"os"
	"os/signal"
	"syscall"
)

var Token string

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	var taskUrl string
	var storyPoints = make(map[string]string)

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		err := handlers.MessageHandler(s, m, &taskUrl, storyPoints)
		if err != nil {
			log.Println(err)
		}
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := handlers.InteractionHandler(s, i, &taskUrl, storyPoints)
		if err != nil {
			log.Println(err)
		}
	})

	err = session.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = session.Close()
	if err != nil {
		return
	}
}
