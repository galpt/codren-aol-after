package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/afero"
	xurls "mvdan.cc/xurls/v2"
)

func katInzYTDL(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".yt") {

		if strings.ToLower(m.Content) == ".yt.help" {
			s.ChannelMessageSendReply(m.ChannelID, "YouTube audio enhancer done right.\n\n**How to Use**\n`.yt <yt link>` — <@!854071193833701416> will enhance the audio in MP3 format;\n\n**Examples**\n`.yt https://youtu.be/qFeKKGDoF2E`\n`.yt https://youtu.be/VfATdDI3604`\n\n**Notes**\n```\n• The process should only takes 10 seconds or less;\n• Files bigger than 8 MB aren't allowed by Discord. Thus, they won't be sent back to you;\n```\n", m.Reference())
		} else {

			if ytLock {
				// if there's a user using the ytdl right now,
				// wait until the process is finished.
				s.ChannelMessageSendReply(m.ChannelID, "There's a user using this feature right now.\nPlease wait until the process is finished.", m.Reference())
			} else {

				ytRelax := xurls.Relaxed()

				ytdlSplit, err := kemoSplit(m.Content, " ")
				if err != nil {
					fmt.Println(" [ytdlSplit] ", err)

					if len(universalLogs) >= universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					ytLock = false
					return
				}

				// rawArgs shouldn't be empty
				if len(ytdlSplit) > 1 {

					if strings.ToLower(ytdlSplit[0]) == ".yt" {

						s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
						ytLock = true

						osFS.RemoveAll("./ytdl")
						osFS.MkdirAll("./ytdl", 0777)
						katInzVidID = ""

						// delete user's message and send confirmation as a reply
						scanLinks := ytRelax.FindAllString(m.Content, -1)

						// get the video ID
						if strings.Contains(scanLinks[0], "www.youtube.com") {
							// sample URL >> https://www.youtube.com/watch?v=J5x0tLiItVY
							splitVidID := strings.Split(scanLinks[0], "youtube.com/watch?v=")
							katInzVidID = splitVidID[1]
						} else if strings.Contains(scanLinks[0], "youtu.be") {
							// sample URL >> https://youtu.be/J5x0tLiItVY
							splitVidID := strings.Split(scanLinks[0], "youtu.be/")
							katInzVidID = splitVidID[1]
						}
						s.ChannelMessageDelete(m.ChannelID, m.ID)
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Processing `%v`. Please wait.", katInzVidID))

						// run the code
						codeExec := time.Now()
						katYT, err := exec.Command("yt-dlp", "--ignore-config", "--no-playlist", "--user-agent", uaChrome, "--max-filesize", "30m", "-P", "./ytdl", "-o", "%(id)s.%(ext)s", "-x", "--audio-format", "mp3", "--audio-quality", "320k", "-N", "10", scanLinks[0]).Output()
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							ytLock = false
							return
						}
						fmt.Println(string(katYT))
						execTime := time.Since(codeExec)

						outIdx, err := afero.ReadDir(osFS, "./ytdl")
						if err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							ytLock = false
							return
						}

						// report after code execution has ended
						// Create the embed templates
						timeElapsedField := discordgo.MessageEmbedField{
							Name:   "Processing Time",
							Value:  fmt.Sprintf("`%v`", execTime),
							Inline: false,
						}
						newsizeField := discordgo.MessageEmbedField{
							Name:   "New Size",
							Value:  fmt.Sprintf("`%v KB | %v MB`", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
							Inline: false,
						}
						fileIDField := discordgo.MessageEmbedField{
							Name:   "File ID",
							Value:  fmt.Sprintf("`%v`", katInzVidID),
							Inline: false,
						}
						linkField := discordgo.MessageEmbedField{
							Name:   "Data in Memory",
							Value:  fmt.Sprintf("https://cdn.castella.network/yt/%v", outIdx[0].Name()),
							Inline: false,
						}
						messageFields := []*discordgo.MessageEmbedField{&timeElapsedField, &newsizeField, &fileIDField, &linkField}

						aoiEmbedFooter := discordgo.MessageEmbedFooter{
							Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
						}

						aoiEmbedAuthor := discordgo.MessageEmbedAuthor{
							URL:     fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
							Name:    fmt.Sprintf("%v#%v", m.Author.Username, m.Author.Discriminator),
							IconURL: fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
						}

						aoiEmbeds := discordgo.MessageEmbed{
							Title:  fmt.Sprintf("%v's YT", botName),
							Color:  0xeb4034,
							Footer: &aoiEmbedFooter,
							Fields: messageFields,
							Author: &aoiEmbedAuthor,
							Image:  &discordgo.MessageEmbedImage{URL: fmt.Sprintf("https://i.ytimg.com/vi_webp/%v/maxresdefault.webp", katInzVidID)},
						}

						s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)

						ytLock = false
					}

				}

			}

		}

	}

}
