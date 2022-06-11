package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Maid-san's handle to auto-check for banned words
func maidsanAutoCheck(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if m.Author.ID == staffID[0] {
		maidsanWatchCurrentUser = maidsanWatchPreviousUser
	} else {
		maidsanWatchCurrentUser = m.Author.Username + "#" + m.Author.Discriminator
		maidsanWatchPreviousUser = m.Author.Username + "#" + m.Author.Discriminator
	}

	// update the bot's status dynamically every 5 seconds
	go updateBotStatus(s, m)

	// Get channel last message IDs
	senderUserID := m.Author.ID
	senderUsername := m.Author.Username + "#" + m.Author.Discriminator
	maidsanLastMsgChannelID = m.ChannelID

	maidsanLastMsgID = m.ID
	maidsanLowercaseLastMsg = strings.ToLower(m.Content)

	scanLinks := xurlsRelaxed.FindAllString(maidsanLowercaseLastMsg, -1)

	katInzBlacklistLinkDetected := false
	for atIdx := range katInzBlacklist {
		for linkIdx := range scanLinks {
			if strings.EqualFold(scanLinks[linkIdx], strings.ToLower(katInzBlacklist[atIdx])) {
				maidsanLowercaseLastMsg = strings.ReplaceAll(maidsanLowercaseLastMsg, katInzBlacklist[atIdx], " [EDITED] ")
				katInzBlacklistLinkDetected = true
			}
		}
	}
	maidsanEditedLastMsg = maidsanLowercaseLastMsg

	if katInzBlacklistLinkDetected {
		// Create the embed templates
		senderField := discordgo.MessageEmbedField{
			Name:   "Sender",
			Value:  fmt.Sprintf("<@%v>", senderUserID),
			Inline: true,
		}
		senderUserIDField := discordgo.MessageEmbedField{
			Name:   "User ID",
			Value:  fmt.Sprintf("%v", senderUserID),
			Inline: true,
		}
		reasonField := discordgo.MessageEmbedField{
			Name:   "Reason",
			Value:  "Blacklisted Links/Banned Words",
			Inline: true,
		}
		editedMsgField := discordgo.MessageEmbedField{
			Name:   "Edited Message",
			Value:  fmt.Sprintf("%v", maidsanEditedLastMsg),
			Inline: false,
		}
		messageFields := []*discordgo.MessageEmbedField{&senderField, &senderUserIDField, &reasonField, &editedMsgField}

		aoiEmbedFooter := discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
		}

		aoiEmbeds := discordgo.MessageEmbed{
			Title:  fmt.Sprintf("Edited by %v ❤️", botName),
			Color:  0x4287f5,
			Footer: &aoiEmbedFooter,
			Fields: messageFields,
		}

		s.ChannelMessageDelete(maidsanLastMsgChannelID, maidsanLastMsgID)
		s.ChannelMessageSendEmbed(maidsanLastMsgChannelID, &aoiEmbeds)

		// Reformat user data before printed out
		userAvatar := m.Author.Avatar
		userisBot := fmt.Sprintf("%v", m.Author.Bot)
		userAccType := ""
		userAvaEmbedImgURL := ""

		// Check whether the user's avatar type is GIF or not
		if strings.Contains(userAvatar, "a_") {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".gif?size=4096"
		} else {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".jpg?size=4096"
		}

		// Check the user's account type
		if userisBot == "true" {
			userAccType = "Bot Account"
		} else {
			userAccType = "Standard User Account"
		}

		// copy logs to Maid-san's memory
		maidsanTranslatedMsg = fmt.Sprintf("https://translate.google.com/?sl=auto&tl=en&text=%v&op=translate", url.QueryEscape(maidsanEditedLastMsg))

		maidsanLogsTemplate = fmt.Sprintf("\n •===========================• \n • Timestamp: %v \n •===========================• \n \n Username: %v \n User ID: %v \n Profile Picture: %v \n Account Type: %v \n Message ID: %v \n Message:\n%v \n Translation:\n%v \n\n", time.Now().UTC().Format(time.RFC850), senderUsername, senderUserID, userAvaEmbedImgURL, userAccType, m.ID, maidsanEditedLastMsg, maidsanTranslatedMsg)

		lastMsgTimestamp = fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850))
		lastMsgUsername = fmt.Sprintf("%v", senderUsername)
		lastMsgUserID = fmt.Sprintf("%v", senderUserID)
		lastMsgpfp = fmt.Sprintf("%v", userAvaEmbedImgURL)
		lastMsgAccType = fmt.Sprintf("%v", userAccType)
		lastMsgID = fmt.Sprintf("%v", m.ID)
		lastMsgContent = fmt.Sprintf("%v", maidsanEditedLastMsg)
		lastMsgTranslation = fmt.Sprintf("%v", maidsanTranslatedMsg)

		if len(maidsanLogs) < maidsanLogsLimit {
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		} else {
			maidsanLogs = nil
			timestampLogs = nil
			useridLogs = nil
			profpicLogs = nil
			acctypeLogs = nil
			msgidLogs = nil
			msgLogs = nil
			translateLogs = nil
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		}
	} else {

		// Reformat user data before printed out
		userAvatar := m.Author.Avatar
		userisBot := fmt.Sprintf("%v", m.Author.Bot)
		userAccType := ""
		userAvaEmbedImgURL := ""

		// Check whether the user's avatar type is GIF or not
		if strings.Contains(userAvatar, "a_") {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".gif?size=4096"
		} else {
			userAvaEmbedImgURL = "https://cdn.discordapp.com/avatars/" + senderUserID + "/" + userAvatar + ".jpg?size=4096"
		}

		// Check the user's account type
		if userisBot == "true" {
			userAccType = "Bot Account"
		} else {
			userAccType = "Standard User Account"
		}

		// copy logs to Maid-san's memory
		maidsanTranslatedMsg = fmt.Sprintf("https://translate.google.com/?sl=auto&tl=en&text=%v&op=translate", url.QueryEscape(m.Content))

		maidsanLogsTemplate = fmt.Sprintf("\n •===========================• \n • Timestamp: %v \n •===========================• \n \n Username: %v \n User ID: %v \n Profile Picture: %v \n Account Type: %v \n Message ID: %v \n Message:\n%v \n Translation:\n%v \n\n", time.Now().UTC().Format(time.RFC850), senderUsername, senderUserID, userAvaEmbedImgURL, userAccType, m.ID, maidsanEditedLastMsg, maidsanTranslatedMsg)

		lastMsgTimestamp = fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC850))
		lastMsgUsername = fmt.Sprintf("%v", senderUsername)
		lastMsgUserID = fmt.Sprintf("%v", senderUserID)
		lastMsgpfp = fmt.Sprintf("%v", userAvaEmbedImgURL)
		lastMsgAccType = fmt.Sprintf("%v", userAccType)
		lastMsgID = fmt.Sprintf("%v", m.ID)
		lastMsgContent = fmt.Sprintf("%v", m.Content)
		lastMsgTranslation = fmt.Sprintf("%v", maidsanTranslatedMsg)

		if len(maidsanLogs) < maidsanLogsLimit {
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		} else {
			maidsanLogs = nil
			timestampLogs = nil
			useridLogs = nil
			profpicLogs = nil
			acctypeLogs = nil
			msgidLogs = nil
			msgLogs = nil
			translateLogs = nil
			maidsanLogs = append(maidsanLogs, maidsanLogsTemplate)
			timestampLogs = append(timestampLogs, lastMsgTimestamp)
			useridLogs = append(useridLogs, lastMsgUserID)
			profpicLogs = append(profpicLogs, lastMsgpfp)
			acctypeLogs = append(acctypeLogs, lastMsgAccType)
			msgidLogs = append(msgidLogs, lastMsgID)
			msgLogs = append(msgLogs, lastMsgContent)
			translateLogs = append(translateLogs, lastMsgTranslation)
		}
	}
}
