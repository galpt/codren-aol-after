package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func updateBotStatus(s *discordgo.Session, m *discordgo.MessageCreate) {

	// prevent spamming the discord API
	time.Sleep(5 * time.Second)

	// dynamically change indicator status
	setActivityText := discordgo.Activity{
		Name: maidsanWatchCurrentUser,
		Type: 3,
	}

	botStatusData := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{&setActivityText},
		Status:     statusSlice[statusInt],
		AFK:        false,
	}
	s.UpdateStatusComplex(botStatusData)

	if statusInt < 2 {
		statusInt++
	} else {
		statusInt = 0
	}

}
