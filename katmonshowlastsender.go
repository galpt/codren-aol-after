package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/afero"
)

func katMonShowLastSender(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".lastsender") {

		if strings.ToLower(m.Content) == ".lastsender" {

			userID := m.Author.ID

			// Only Creator-sama who has the permission
			if strings.Contains(userID, staffID[0]) {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

				osFS.RemoveAll("./logs")
				osFS.MkdirAll("./logs", 0777)

				// ==================================
				// Create a new logs.txt
				createLogsFile, err := osFS.Create("./logs/logs.txt")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) >= universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				} else {
					// Write to the file
					writeLogsFile, err := createLogsFile.WriteString(fmt.Sprintf("%v", maidsanLogs))
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) >= universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}

						return
					} else {
						// Close the file
						if err := createLogsFile.Close(); err != nil {
							fmt.Println(" [ERROR] ", err)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
							}

							return
						} else {
							winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", createLogsFile.Name(), (writeLogsFile / Kilobyte), (writeLogsFile / Megabyte))
							fmt.Println(winLogs)

							if len(universalLogs) >= universalLogsLimit {
								universalLogs = nil
							} else {
								universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
							}
						}
					}
				}

				outIdx, err := afero.ReadDir(osFS, "./logs")
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) >= universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}

				readOutput, err := afero.ReadFile(osFS, fmt.Sprintf("./logs/%v", outIdx[0].Name()))
				if err != nil {
					fmt.Println(" [ERROR] ", err)

					if len(universalLogs) >= universalLogsLimit {
						universalLogs = nil
					} else {
						universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
					}

					return
				}
				reader := bytes.NewReader(readOutput)

				// add some checks to prevent panics.
				// panic: runtime error: index out of range [-2]
				if len(maidsanLogs) >= 2 {

					// report after code execution has ended
					// Create the embed templates
					usernameField := discordgo.MessageEmbedField{
						Name:   "Data Issuer",
						Value:  fmt.Sprintf("<@!%v>", userID),
						Inline: false,
					}
					lastsenderField := discordgo.MessageEmbedField{
						Name:   "Last Sender",
						Value:  fmt.Sprintf("<@!%v>", useridLogs[(len(useridLogs)-2)]),
						Inline: false,
					}
					timestampField := discordgo.MessageEmbedField{
						Name:   "Timestamp",
						Value:  fmt.Sprintf("`%v`", timestampLogs[(len(timestampLogs)-2)]),
						Inline: false,
					}
					pfpField := discordgo.MessageEmbedField{
						Name:   "Profile Picture",
						Value:  fmt.Sprintf("```\n%v\n```", profpicLogs[(len(profpicLogs)-2)]),
						Inline: false,
					}
					acctypeField := discordgo.MessageEmbedField{
						Name:   "Account Type",
						Value:  fmt.Sprintf("`%v`", acctypeLogs[(len(acctypeLogs)-2)]),
						Inline: false,
					}
					msgidField := discordgo.MessageEmbedField{
						Name:   "Message ID",
						Value:  fmt.Sprintf("`%v`", msgidLogs[(len(msgidLogs)-2)]),
						Inline: false,
					}
					msgcontentField := discordgo.MessageEmbedField{
						Name:   "Message",
						Value:  fmt.Sprintf("```\n%v\n```", msgLogs[(len(msgLogs)-2)]),
						Inline: false,
					}
					translateField := discordgo.MessageEmbedField{
						Name:   "Translation",
						Value:  fmt.Sprintf("```\n%v\n```", translateLogs[(len(translateLogs)-2)]),
						Inline: false,
					}
					logsindexField := discordgo.MessageEmbedField{
						Name:   "Logs Limit",
						Value:  fmt.Sprintf("`%v / %v`", len(maidsanLogs), maidsanLogsLimit),
						Inline: false,
					}
					logssizeField := discordgo.MessageEmbedField{
						Name:   "Logs Size",
						Value:  fmt.Sprintf("`%v KB | %v MB`", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
						Inline: false,
					}
					messageFields := []*discordgo.MessageEmbedField{&usernameField, &lastsenderField, &timestampField, &pfpField, &acctypeField, &msgidField, &msgcontentField, &translateField, &logsindexField, &logssizeField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  fmt.Sprintf("All Seeing Eyes of %v", botName),
						Color:  0x4287f5,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
					s.ChannelFileSend(m.ChannelID, outIdx[0].Name(), reader)
				} else if len(maidsanLogs) >= 1 {

					// report after code execution has ended
					// Create the embed templates
					usernameField := discordgo.MessageEmbedField{
						Name:   "Data Issuer",
						Value:  fmt.Sprintf("<@!%v>", userID),
						Inline: false,
					}
					lastsenderField := discordgo.MessageEmbedField{
						Name:   "Last Sender",
						Value:  fmt.Sprintf("<@!%v>", useridLogs[(len(useridLogs)-1)]),
						Inline: false,
					}
					timestampField := discordgo.MessageEmbedField{
						Name:   "Timestamp",
						Value:  fmt.Sprintf("`%v`", timestampLogs[(len(timestampLogs)-1)]),
						Inline: false,
					}
					pfpField := discordgo.MessageEmbedField{
						Name:   "Profile Picture",
						Value:  fmt.Sprintf("```\n%v\n```", profpicLogs[(len(profpicLogs)-1)]),
						Inline: false,
					}
					acctypeField := discordgo.MessageEmbedField{
						Name:   "Account Type",
						Value:  fmt.Sprintf("`%v`", acctypeLogs[(len(acctypeLogs)-1)]),
						Inline: false,
					}
					msgidField := discordgo.MessageEmbedField{
						Name:   "Message ID",
						Value:  fmt.Sprintf("`%v`", msgidLogs[(len(msgidLogs)-1)]),
						Inline: false,
					}
					msgcontentField := discordgo.MessageEmbedField{
						Name:   "Message",
						Value:  fmt.Sprintf("```\n%v\n```", msgLogs[(len(msgLogs)-1)]),
						Inline: false,
					}
					translateField := discordgo.MessageEmbedField{
						Name:   "Translation",
						Value:  fmt.Sprintf("```\n%v\n```", translateLogs[(len(translateLogs)-1)]),
						Inline: false,
					}
					logsindexField := discordgo.MessageEmbedField{
						Name:   "Logs Limit",
						Value:  fmt.Sprintf("`%v / %v`", len(maidsanLogs), maidsanLogsLimit),
						Inline: false,
					}
					logssizeField := discordgo.MessageEmbedField{
						Name:   "Logs Size",
						Value:  fmt.Sprintf("`%v KB | %v MB`", (outIdx[0].Size() / Kilobyte), (outIdx[0].Size() / Megabyte)),
						Inline: false,
					}
					messageFields := []*discordgo.MessageEmbedField{&usernameField, &lastsenderField, &timestampField, &pfpField, &acctypeField, &msgidField, &msgcontentField, &translateField, &logsindexField, &logssizeField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  fmt.Sprintf("All Seeing Eyes of %v", botName),
						Color:  0x4287f5,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
					s.ChannelFileSend(m.ChannelID, outIdx[0].Name(), reader)
				} else {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I couldn't get any data from my memory.\n```\nLogs Data: %v / %v\n```", len(maidsanLogs), maidsanLogsLimit))
				}
			}
		}
	}

}
