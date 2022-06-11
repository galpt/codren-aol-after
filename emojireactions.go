package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// react with the available server emojis
func emojiReactions(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else {

		// React with ganyustare emoji
		// if the m.Content contains "geez" word
		if strings.Contains(strings.ToLower(m.Content), "geez") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "ganyustare:903098908966785024")
		} else if strings.Contains(strings.ToLower(m.Content), "<:ganyustare:903098908966785024>") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "ganyustare:903098908966785024")
		}
	}

}
