package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Maid-san's emoji reactions handler
func maidsanEmojiReact(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else {

		// reply with custom castella network's sticker list
		if strings.Contains(m.Content, ".stk") {
			tempmsg1 := strings.ReplaceAll(m.Content, ".stk", "")
			tempmsg2 := strings.ReplaceAll(tempmsg1, " ", "")
			realMsg := tempmsg2

			for stkIdx := range stickerList {
				if strings.Contains(stickerList[stkIdx], realMsg) {
					s.ChannelMessageSendReply(m.ChannelID, stickerList[stkIdx], m.Reference())
					break
				}
			}
		}

		customEmojiDetected = false

		// Reply with custom emoji if the message contains the keyword
		for currIdx := range maidsanEmojiInfo {
			replyremoveNewLines = strings.ReplaceAll(maidsanEmojiInfo[currIdx], "\n", "")
			replyremoveSpaces = strings.ReplaceAll(replyremoveNewLines, " ", "")
			replysplitEmojiInfo = strings.Split(replyremoveSpaces, "——")

			if strings.EqualFold(replysplitEmojiInfo[0], strings.ToLower(m.Content)) {
				customEmojiDetected = true
				if replysplitEmojiInfo[2] != "false" {
					customEmojiReply = fmt.Sprintf("<a:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1])
				} else {
					customEmojiReply = fmt.Sprintf("<:%v:%v>", replysplitEmojiInfo[0], replysplitEmojiInfo[1])
				}
			}
		}

		if customEmojiDetected {
			s.ChannelMessageSend(m.ChannelID, customEmojiReply)
		} else {
			s.MessageReactionAdd(m.ChannelID, m.ID, customEmojiSlice[customEmojiIdx])
			if customEmojiIdx == (len(customEmojiSlice) - 1) {
				customEmojiIdx = 0
			} else {
				customEmojiIdx++
			}
		}

	}

}
