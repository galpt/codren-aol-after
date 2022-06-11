package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/afero"
)

// support for xvid feature
func xvid(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".xv") {

		var (
			userID      = m.Author.ID
			xvURL       = ""
			xvTotalSize = ""
			xvVidName   = ""
		)

		xvsplitText, err := kemoSplit(m.Content, " ")
		if err != nil {
			fmt.Println(" [xvsplitText] ", err)

			if len(universalLogs) >= universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}

			return
		}

		// rawArgs shouldn't be empty
		if len(xvsplitText) > 1 {

			if strings.Contains(strings.ToLower(xvsplitText[0]), ".xv") {

				if strings.Contains(strings.ToLower(xvsplitText[1]), "help") {

					s.ChannelMessageSendReply(m.ChannelID, "**XV**\nAn accelerator for `xvideos.com` done right.\n\n**How to Use**\n`.xv help` — show the help message;\n`.xv deldata` — delete cached data;\n`.xv <video URL>` — reconstruct video & cache it to boost performance;\n\n**Example**\n`.xv https://www.xvideos.com/video54147993/sexy_solo_babe_masturbating` — <@!854071193833701416> will reconstruct the video for you;\n", m.Reference())

				} else if strings.Contains(strings.ToLower(xvsplitText[1]), "deldata") {

					if xvLock {
						// if there's a user using the NH right now,
						// wait until the process is finished.
						s.ChannelMessageSendReply(m.ChannelID, "There's a user using this feature right now.\nPlease wait until the process is finished.", m.Reference())
					} else {

						if strings.Contains(userID, staffID[0]) {
							s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

							var (
								cacheList []string
							)

							// start counting time elapsed
							codeExec := time.Now()

							// read cache dir
							readDir1, err := afero.ReadDir(osFS, "./xvids/")
							if err != nil {
								fmt.Println(" [readDir1] ", err)

								if len(universalLogs) >= universalLogsLimit {
									universalLogs = nil
								} else {
									universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
								}

								return
							}

							for idx := range readDir1 {
								cacheList = append(cacheList, fmt.Sprintf("%v\n", readDir1[idx].Name()))
							}

							osFS.RemoveAll("./xvids/")
							osFS.MkdirAll("./xvids/", 0777)

							// Create the embed templates.
							timeElapsedField := discordgo.MessageEmbedField{
								Name:   "Processing Time",
								Value:  fmt.Sprintf("`%v`", time.Since(codeExec)),
								Inline: false,
							}
							cacheSizeField := discordgo.MessageEmbedField{
								Name:   "Total Cache",
								Value:  fmt.Sprintf("`%v Cache(s)`", len(readDir1)),
								Inline: false,
							}
							messageFields := []*discordgo.MessageEmbedField{&timeElapsedField, &cacheSizeField}

							aoiEmbedFooter := discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
							}

							aoiEmbedAuthor := discordgo.MessageEmbedAuthor{
								URL:     fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
								Name:    fmt.Sprintf("%v#%v", m.Author.Username, m.Author.Discriminator),
								IconURL: fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
							}

							aoiEmbeds := discordgo.MessageEmbed{
								Title:  "XV",
								Color:  0x82ff86,
								Footer: &aoiEmbedFooter,
								Fields: messageFields,
								Author: &aoiEmbedAuthor,
							}

							s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
							s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("**Deleted Cache(s)**\n```\n%v\n```", cacheList), m.Reference())

						} else {
							// only for Creator-sama
							s.ChannelMessageSendReply(m.ChannelID, "You are not allowed to access this command.", m.Reference())
						}

					}

				} else if strings.Contains(strings.ToLower(xvsplitText[1]), "https://") {

					s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
					xvURL = xvsplitText[1]

					if xvLock {
						// if there's a user using the NH right now,
						// wait until the process is finished.
						s.ChannelMessageSendReply(m.ChannelID, "There's a user using this feature right now.\nPlease wait until the process is finished.", m.Reference())
					} else {

						// lock to prevent race condition
						xvLock = true

						// start counting time elapsed
						codeExec := time.Now()

						// send a quick message reply as a confirmation
						s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Fetching `%v` data.\nMaybe you can make a cup of tea while I'm working on it.", xvURL), m.Reference())

						// make a new folder
						osFS.RemoveAll(fmt.Sprintf("./xvids/%v/", userID))
						osFS.MkdirAll(fmt.Sprintf("./xvids/%v/", userID), 0777)

						// run the code
						katXV := exec.Command("yt-dlp", "--ignore-config", "--no-playlist", "--user-agent", uaChrome, "-P", fmt.Sprintf("./xvids/%v", userID), "-o", "%(duration)s.%(filesize)s.%(resolution)s.%(id)s.%(ext)s", "-N", "10", "-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", xvURL)
						output, err := katXV.CombinedOutput()
						if err != nil {
							errMsg := fmt.Sprintf(" [katXV] %v: %v", err, string(output))
							fmt.Println(errMsg)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						}
						fmt.Println(string(output))

						chkFile, err := afero.ReadDir(osFS, fmt.Sprintf("./xvids/%v", userID))
						if err != nil {
							fmt.Println(" [chkFile] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							xvLock = false
							return
						}
						xvVidName = chkFile[0].Name()
						xvTotalSize = fmt.Sprintf("%v KB | %v MB", (chkFile[0].Size() / Kilobyte), (chkFile[0].Size() / Megabyte))

						// get time elapsed data
						execTime := time.Since(codeExec)

						// Create the embed templates.
						timeElapsedField := discordgo.MessageEmbedField{
							Name:   "Processing Time",
							Value:  fmt.Sprintf("`%v`", execTime),
							Inline: false,
						}
						sizeField := discordgo.MessageEmbedField{
							Name:   "Total Size",
							Value:  fmt.Sprintf("`%v`", xvTotalSize),
							Inline: false,
						}
						urlField := discordgo.MessageEmbedField{
							Name:   "Data in Memory",
							Value:  fmt.Sprintf("https://cdn.castella.network/xv/%v/%v", userID, xvVidName),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&timeElapsedField, &sizeField, &urlField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbedAuthor := discordgo.MessageEmbedAuthor{
							URL:     fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
							Name:    fmt.Sprintf("%v#%v", m.Author.Username, m.Author.Discriminator),
							IconURL: fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  "XV",
							Color:  0xf06967,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
							Author: &aoiEmbedAuthor,
						}

						s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)

						// unlock after the process is finished
						xvLock = false

					}

				}

			}

		}
	}

}
