package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
)

// Get current COVID-19 data for Indonesia country
func getCovidData(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".covid19") {

		getcovidinputSplit, err := kemoSplit(m.Content, " ")
		if err != nil {
			fmt.Println(" [getcovidinputSplit] ", err)

			if len(universalLogs) >= universalLogsLimit {
				universalLogs = nil
			} else {
				universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
			}

			return
		}

		if strings.ToLower(getcovidinputSplit[0]) == ".covid19" {

			// countryArgs shouldn't be empty
			if len(getcovidinputSplit) > 1 {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

				countryArgs := getcovidinputSplit[1]

				if countryArgs == "indonesia" {
					// Get covid-19 json data Indonesia
					covIndo, err := httpclient.Get("https://data.covid19.go.id/public/api/update.json")
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) >= universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					}

					bodyCovIndo, err := ioutil.ReadAll(covIndo.Body)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) >= universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					}

					// Indonesia - Reformat JSON before printed out
					indoCreatedVal := gjson.Get(string(bodyCovIndo), `update.penambahan.created`)
					indoPosVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_positif`)
					indoMeninggalVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_meninggal`)
					indoSembuhVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_sembuh`)
					indoDirawatVal := gjson.Get(string(bodyCovIndo), `update.penambahan.jumlah_dirawat`)
					indoTotalPosVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_positif`)
					indoTotalMeninggalVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_meninggal`)
					indoTotalSembuhVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_sembuh`)
					indoTotalDirawatVal := gjson.Get(string(bodyCovIndo), `update.total.jumlah_dirawat`)

					// Create the embed templates
					createdField := discordgo.MessageEmbedField{
						Name:   "Date Created",
						Value:  indoCreatedVal.String(),
						Inline: true,
					}
					countryField := discordgo.MessageEmbedField{
						Name:   "Country",
						Value:  strings.ToUpper(countryArgs),
						Inline: true,
					}
					totalConfirmedField := discordgo.MessageEmbedField{
						Name:   "Total Confirmed",
						Value:  fmt.Sprintf("%v", indoTotalPosVal.Int()),
						Inline: true,
					}
					totalDeathsField := discordgo.MessageEmbedField{
						Name:   "Total Deaths",
						Value:  fmt.Sprintf("%v", indoTotalMeninggalVal.Int()),
						Inline: true,
					}
					totalRecoveredField := discordgo.MessageEmbedField{
						Name:   "Total Recovered",
						Value:  fmt.Sprintf("%v", indoTotalSembuhVal.Int()),
						Inline: true,
					}
					totalTreatedField := discordgo.MessageEmbedField{
						Name:   "Total Treated",
						Value:  fmt.Sprintf("%v", indoTotalDirawatVal.Int()),
						Inline: true,
					}
					additionalConfirmedField := discordgo.MessageEmbedField{
						Name:   "Additional Confirmed",
						Value:  fmt.Sprintf("%v", indoPosVal.Int()),
						Inline: true,
					}
					additionalDeathsField := discordgo.MessageEmbedField{
						Name:   "Additional Deaths",
						Value:  fmt.Sprintf("%v", indoMeninggalVal.Int()),
						Inline: true,
					}
					additionalRecoveredField := discordgo.MessageEmbedField{
						Name:   "Additional Recovered",
						Value:  fmt.Sprintf("%v", indoSembuhVal.Int()),
						Inline: true,
					}
					additionalTreatedField := discordgo.MessageEmbedField{
						Name:   "Additional Treated",
						Value:  fmt.Sprintf("%v", indoDirawatVal.Int()),
						Inline: true,
					}
					messageFields := []*discordgo.MessageEmbedField{&createdField, &countryField, &totalConfirmedField, &totalDeathsField, &totalRecoveredField, &totalTreatedField, &additionalConfirmedField, &additionalDeathsField, &additionalRecoveredField, &additionalTreatedField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  "Latest COVID-19 Data",
						Color:  0xE06666,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
					covIndo.Body.Close()
				} else {
					// Get covid-19 json data from a certain country
					// based on the user's argument
					urlCountry := "https://covid19.mathdro.id/api/countries/" + countryArgs
					covData, err := httpclient.Get(urlCountry)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) >= universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					}

					bodyCovData, err := ioutil.ReadAll(covData.Body)
					if err != nil {
						fmt.Println(" [ERROR] ", err)

						if len(universalLogs) >= universalLogsLimit {
							universalLogs = nil
						} else {
							universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
						}
					}

					// Reformat JSON before printed out
					countryCreatedVal := gjson.Get(string(bodyCovData), `lastUpdate`)
					countryTotalPosVal := gjson.Get(string(bodyCovData), `confirmed.value`)
					countryTotalSembuhVal := gjson.Get(string(bodyCovData), `recovered.value`)
					countryTotalMeninggalVal := gjson.Get(string(bodyCovData), `deaths.value`)

					// Create the embed templates
					createdField := discordgo.MessageEmbedField{
						Name:   "Date Created",
						Value:  countryCreatedVal.String(),
						Inline: true,
					}
					countryField := discordgo.MessageEmbedField{
						Name:   "Country",
						Value:  strings.ToUpper(countryArgs),
						Inline: true,
					}
					totalConfirmedField := discordgo.MessageEmbedField{
						Name:   "Total Confirmed",
						Value:  fmt.Sprintf("%v", countryTotalPosVal.Int()),
						Inline: true,
					}
					totalDeathsField := discordgo.MessageEmbedField{
						Name:   "Total Deaths",
						Value:  fmt.Sprintf("%v", countryTotalMeninggalVal.Int()),
						Inline: true,
					}
					totalRecoveredField := discordgo.MessageEmbedField{
						Name:   "Total Recovered",
						Value:  fmt.Sprintf("%v", countryTotalSembuhVal.Int()),
						Inline: true,
					}
					messageFields := []*discordgo.MessageEmbedField{&createdField, &countryField, &totalConfirmedField, &totalDeathsField, &totalRecoveredField}

					aoiEmbedFooter := discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
					}

					aoiEmbeds := discordgo.MessageEmbed{
						Title:  "Latest COVID-19 Data",
						Color:  0xE06666,
						Footer: &aoiEmbedFooter,
						Fields: messageFields,
					}

					s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)
					covData.Body.Close()
				}

			}
		}

	}

}
