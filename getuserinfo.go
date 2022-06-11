package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Get a brief information about the mentioned user
func getUserInfo(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".check") {

		splitgetUserInfo, err := kemoSplit(m.Content, " ")
		if err != nil {
			fmt.Println(" [splitgetUserInfo] ", err)
			if len(universalLogs) >= universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}
			return
		}

		// rawArgs shouldn't be empty
		if len(splitgetUserInfo) > 1 {
			if strings.ToLower(splitgetUserInfo[0]) == ".check" {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

				var finalUID = ""
				getUID := re.FindAllString(splitgetUserInfo[1], -1)

				for idx := range getUID {
					finalUID += getUID[idx]
				}

				userData, err := s.User(finalUID)
				if err != nil {
					fmt.Println(" [userData] ", err)
					if len(universalLogs) >= universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}
					return
				}

				// Reformat user data before printed out
				userUsername := userData.Username + "#" + userData.Discriminator
				userID := userData.ID
				userAvatar := userData.Avatar
				userisBot := fmt.Sprintf("%v", userData.Bot)
				userAccType := ""
				userAvatarURLFullSize := ""
				userAvaEmbedImgURL := ""

				// Check whether the user's avatar type is GIF or not
				if strings.Contains(userAvatar, "a_") {
					userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=4096"
					userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".gif?size=256"
				} else {
					userAvatarURLFullSize = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=4096"
					userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + userID + "/" + userAvatar + ".jpg?size=256"
				}

				// Check the user's account type
				if userisBot == "true" {
					userAccType = "Bot Account"
				} else {
					userAccType = "Standard User Account"
				}

				// Create the embed templates
				usernameField := discordgo.MessageEmbedField{
					Name:   "Username",
					Value:  userUsername,
					Inline: true,
				}
				userIDField := discordgo.MessageEmbedField{
					Name:   "User ID",
					Value:  userID,
					Inline: true,
				}
				userAvatarField := discordgo.MessageEmbedField{
					Name:   "Profile Picture URL",
					Value:  userAvatarURLFullSize,
					Inline: false,
				}
				userAccTypeField := discordgo.MessageEmbedField{
					Name:   "Account Type",
					Value:  userAccType,
					Inline: true,
				}
				messageFields := []*discordgo.MessageEmbedField{&usernameField, &userIDField, &userAvatarField, &userAccTypeField}

				aoiEmbedFooter := discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
				}

				aoiEmbedThumbnail := discordgo.MessageEmbedThumbnail{
					URL: userAvaEmbedImgURL,
				}

				aoiEmbeds := discordgo.MessageEmbed{
					Title:     "About User",
					Color:     0x00D2FF,
					Thumbnail: &aoiEmbedThumbnail,
					Footer:    &aoiEmbedFooter,
					Fields:    messageFields,
				}

				s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
			}

		}

	}

}
