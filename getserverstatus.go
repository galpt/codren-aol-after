package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	totalmem "github.com/pbnjay/memory"
	"github.com/showwin/speedtest-go/speedtest"
)

// Get realtime server status
func getServerStatus(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".status") {

		user, _ := speedtest.FetchUserInfo()

		if svstatLock {
			// if there's a user using this feature right now,
			// wait until the process is finished.
			s.ChannelMessageSendReply(m.ChannelID, "There's a user using this feature right now.\nPlease wait until the process is finished.", m.Reference())
		} else {
			svstatLock = true

			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

			serverList, _ := speedtest.FetchServers(user)
			targets, _ := serverList.FindServer([]int{})
			var speedResult string

			for _, s := range targets {
				s.PingTest()
				s.DownloadTest(false)
				s.UploadTest(false)

				speedResult = fmt.Sprintf("Latency: %s\nDownload: %.1f Mbps\nUpload: %.1f Mbps\n", s.Latency, s.DLSpeed, s.ULSpeed)
			}

			runtime.ReadMemStats(&mem)
			timeSince := time.Since(duration)

			// Create the embed templates
			cpuCoresField := discordgo.MessageEmbedField{
				Name:   "Available CPU Cores",
				Value:  fmt.Sprintf("`%v`", runtime.NumCPU()),
				Inline: false,
			}
			osMemoryField := discordgo.MessageEmbedField{
				Name:   "Available OS Memory",
				Value:  fmt.Sprintf("`%v MB | %v GB`", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte)),
				Inline: false,
			}
			timeElapsedField := discordgo.MessageEmbedField{
				Name:   "Time Elapsed",
				Value:  fmt.Sprintf("`%v`", timeSince),
				Inline: false,
			}
			netSpeed := discordgo.MessageEmbedField{
				Name:   "Internet Speed",
				Value:  fmt.Sprintf("```\n%v\n```", speedResult),
				Inline: false,
			}
			messageFields := []*discordgo.MessageEmbedField{&cpuCoresField, &osMemoryField, &timeElapsedField, &netSpeed}

			aoiEmbedFooter := discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
			}

			aoiEmbeds := discordgo.MessageEmbed{
				Title:  fmt.Sprintf("%v's Reports", botName),
				Color:  0xF6B26B,
				Footer: &aoiEmbedFooter,
				Fields: messageFields,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)

			svstatLock = false
		}

	}

}
