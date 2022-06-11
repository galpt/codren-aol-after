package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time one
// of our shards connects.
func onConnect(s *discordgo.Session, evt *discordgo.Connect) {
	fmt.Printf("[INFO] Shard #%v connected.\n", s.ShardID)

	// reconnect websocket on errors and some other tweaks
	s.ShouldReconnectOnError = true
	s.Identify.Compress = true
	s.Identify.Properties.Browser = "Discord iOS"

	if len(universalLogs) >= universalLogsLimit {
		universalLogs = nil
	} else {
		universalLogs = append(universalLogs, fmt.Sprintf("\n[INFO] Shard #%v connected. | Connected shards: %v", s.ShardID, s.ShardCount))
	}

	setActivityText := discordgo.Activity{
		Name: fmt.Sprintf("%v thread(s)", s.ShardCount),
		Type: 3,
	}

	botStatusData := discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{&setActivityText},
		Status:     "dnd",
		AFK:        false,
	}
	s.UpdateStatusComplex(botStatusData)

	// autocheck all emojis from the guilds the bot is in
	// clear slices
	maidsanEmojiInfo = nil
	customEmojiSlice = nil

	// get guild list
	getGuilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println(" [getGuilds] ", err)

		if len(universalLogs) >= universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	for guildIdx := range getGuilds {

		// Check the available emoji list
		getEmoji, err := s.GuildEmojis(getGuilds[guildIdx].ID)
		if err != nil {
			fmt.Println(" [getEmoji] ", err)

			if len(universalLogs) >= universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}

			return
		}

		for idxEmoji := range getEmoji {
			maidsanEmojiInfo = append(maidsanEmojiInfo, fmt.Sprintf("\n%v —— %v —— %v —— %v —— %v", getEmoji[idxEmoji].Name, getEmoji[idxEmoji].ID, getEmoji[idxEmoji].Animated, getGuilds[guildIdx].Name, getGuilds[guildIdx].ID))

			customEmojiSlice = append(customEmojiSlice, fmt.Sprintf("%v:%v", getEmoji[idxEmoji].Name, getEmoji[idxEmoji].ID))
		}
	}

	// set custom http client and user agent
	s.Client = httpclient
	s.UserAgent = uaChrome

}
