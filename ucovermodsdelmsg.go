package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	xurls "mvdan.cc/xurls/v2"
)

// Undercover Mods are allowed to delete inappropriate messages.
func ucoverModsDelMsg(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".delmsg") {

		delmsgRelax := xurls.Relaxed()
		userID := m.Author.ID

		delmsginputSplit, err := kemoSplit(m.Content, " ")
		if err != nil {
			fmt.Println(" [delmsginputSplit] ", err)

			if len(universalLogs) >= universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}

			return
		}

		maidsanWatchCurrentUser = "@everyone" // to keep undermods hidden

		if len(delmsginputSplit) > 1 {
			if strings.ToLower(delmsginputSplit[0]) == ".delmsg" {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

				// Check userID in ucoverNewAdded slice
				for chkIdx := range staffID {
					if userID == staffID[0] {

						s.ChannelMessageDelete(m.ChannelID, m.ID)

						scanLinks := delmsgRelax.FindAllString(m.Content, -1)
						splitData, err := kemoSplit(scanLinks[0], "/")
						if err != nil {
							fmt.Println(" [splitData] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						}

						s.ChannelMessageDelete(splitData[len(splitData)-2], splitData[len(splitData)-1])
						maidsanBanUserMsg = fmt.Sprintf("I've deleted MessageID `%v` from <#%v>, Master.", splitData[len(splitData)-1], splitData[len(splitData)-2])
						s.ChannelMessageSend(m.ChannelID, maidsanBanUserMsg)

						// Create the embed templates
						usernameField := discordgo.MessageEmbedField{
							Name:   "Username",
							Value:  fmt.Sprintf("<@!%v>", userID),
							Inline: false,
						}
						modIDField := discordgo.MessageEmbedField{
							Name:   "Undercover ID",
							Value:  fmt.Sprintf("U-%v", chkIdx),
							Inline: false,
						}
						delmsgIDField := discordgo.MessageEmbedField{
							Name:   "Deleted Message ID",
							Value:  fmt.Sprintf("`%v`", splitData[len(splitData)-1]),
							Inline: false,
						}
						delmsgChanField := discordgo.MessageEmbedField{
							Name:   "Deleted Message Channel",
							Value:  fmt.Sprintf("<#%v>", splitData[len(splitData)-2]),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&usernameField, &modIDField, &delmsgIDField, &delmsgChanField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "Usage Information",
							Color:  0x32a852,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
						}

						// Send notification to galpt.
						// We create the private channel with the user who sent the message.
						channel, err := s.UserChannelCreate(staffID[0])
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
							return
						}
						// Then we send the message through the channel we created.
						_, err = s.ChannelMessageSendEmbed(channel.ID, &aoiEmbeds)
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}
						}

						break
					}
				}
			}

		}

	}

}
