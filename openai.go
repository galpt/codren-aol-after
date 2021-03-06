package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rivo/uniseg"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/spf13/afero"
)

// support for OpenAI GPT-3 API
func openAI(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.Contains(m.Content, ".ask") {

		if m.Content == ".ask.help" {

			s.ChannelMessageSendReply(m.ChannelID, "An AI for <@!854071193833701416> done right.\n\n**How to Use**\n`.ask anything` — <@!854071193833701416> will try to answer your request smartly;\n`.ask.clem anything` — <@!854071193833701416> will try to answer in clever mode;\n`.ask.crem anything` — <@!854071193833701416> will try to answer in creative mode;\n`.ask.code.fast anything` — <@!854071193833701416> will try to generate the code faster at the cost of lower answer quality;\n`.ask.code.best anything` — <@!854071193833701416> will try to generate the code better at the cost of slower processing time;\n\n**Examples (General)**\n`.ask How big is Google?`\n`.ask Write a story about a girl named Castella.`\n```css\n.ask Translate this to Japanese:\n\n---\nGood morning!\n---\n\n```\n**Examples (Code Generation)**\n```css\n.ask.code.fast Write a piece of code in Java programming language:\n\n---\nPrint 'Hello, Castella!' to the user using for loop 5 times.\n---\n```\n```css\n.ask.code.fast\n\n---\nTable customers, columns = [CustomerId, FirstName, LastName]\nCreate a MySQL query for a customer named Castella.\n---\nquery =\n```\n**Notes**\n```\n• Answers are 100% generated by AI and might not be accurate;\n• Answers may vary depending on the given clues;\n• Requests submitted may be used to train and improve future models;\n• Most models' training data cuts off in October 2019, so they may not have knowledge of current events.\n```\n", m.Reference())

		} else if strings.Contains(m.Content, ".ask") {

			userID := m.Author.ID
			msgAttachment := m.Attachments

			memFS.RemoveAll("./OpenAI")
			memFS.MkdirAll("./OpenAI", 0777)

			openAIinputSplit, err := kemoSplit(m.Content, " ")
			if err != nil {
				fmt.Println(" [openAIinputSplit] ", err)

				if len(universalLogs) >= universalLogsLimit {
					universalLogs = nil
				} else {
					universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
				}

				return
			}

			if strings.Contains(openAIinputSplit[0], ".ask") {

				var (
					apiKey        = "" // your api key here
					usrInput      = ""
					model         = ""
					mode          = "balanced"
					respEdited    = ""
					allowedTokens = 250 // according to OpenAI's usage guidelines
					charCount     = 0
					costCount     = 0.0
					nvalptr       = 1
					tempptr       = float32(0.3)
					toppptr       = float32(1)
					isCodex       = false
					sendCodex     = false
					//wordCount     = 0
				)

				if strings.Contains(strings.ToLower(m.Content), ".ask.add") {

					if strings.Contains(userID, staffID[0]) {

						var finalUID = ""
						getUID := re.FindAllString(openAIinputSplit[1], -1)

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

						// check if finalUID does exist in openAIAccess
						for chk := range openAIAccess {
							if strings.Contains(finalUID, openAIAccess[chk]) {
								s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<:ganyustare:903098908966785024> `%v#%v` is already allowed to access my knowledge.", userData.Username, userData.Discriminator), m.Reference())
								return
							}
						}

						openAIAccess = append(openAIAccess, finalUID)

						s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<:ganyustare:903098908966785024> I've allowed `%v#%v` to access my knowledge.", userData.Username, userData.Discriminator), m.Reference())
						return
					} else {
						return
					}

				} else if strings.Contains(strings.ToLower(m.Content), ".ask.del") {

					if strings.Contains(userID, staffID[0]) {

						var (
							finalUID = ""
							newUID   []string
						)
						getUID := re.FindAllString(openAIinputSplit[1], -1)

						for idx := range getUID {
							finalUID += getUID[idx]
						}

						for idIDX := range openAIAccess {
							if strings.Contains(finalUID, openAIAccess[idIDX]) {
								newUID = RemoveIndex(openAIAccess, idIDX)
								break
							}
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

						openAIAccess = nil
						openAIAccess = append(openAIAccess, newUID...)
						newUID = nil

						s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<:ganyustare:903098908966785024> Okay, now `%v#%v` won't be able to access my knowledge.", userData.Username, userData.Discriminator), m.Reference())
						return
					} else {
						return
					}

				} else if strings.Contains(strings.ToLower(m.Content), ".ask.clem") {

					mode = "clever"
					tempptr = float32(0.1)
					usrInput = strings.ReplaceAll(m.Content, ".ask.clem", "")

				} else if strings.Contains(strings.ToLower(m.Content), ".ask.crem") {

					mode = "creative"
					tempptr = float32(0.9)
					usrInput = strings.ReplaceAll(m.Content, ".ask.crem", "")

				} else if strings.Contains(strings.ToLower(m.Content), ".ask.code.fast") {

					isCodex = true
					sendCodex = true
					mode = "code-fast"
					model = "code-cushman-001"
					tempptr = float32(0.0)
					allowedTokens = 1000 // according to OpenAI's usage guidelines
					usrInput = strings.ReplaceAll(m.Content, ".ask.code.fast", "")

				} else if strings.Contains(strings.ToLower(m.Content), ".ask.code.best") {

					isCodex = true
					sendCodex = true
					mode = "code-best"
					model = "code-davinci-002"
					tempptr = float32(0.0)
					allowedTokens = 1000 // according to OpenAI's usage guidelines
					usrInput = strings.ReplaceAll(m.Content, ".ask.code.best", "")

				} else {

					mode = "balanced"
					tempptr = float32(0.7)
					usrInput = strings.ReplaceAll(m.Content, ".ask", "")

				}

				// Only Creator-sama who has the permission
				for idx := range openAIAccess {

					if strings.Contains(openAIAccess[idx], userID) {
						s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

						// add check for max current requests in queue
						if currReqPerMin < maxReqPerMin {

							// increase the counter to limit next request
							currReqPerMin = currReqPerMin + 1

							// start counting time elapsed
							codeExec := time.Now()

							// input request shouldn't be more than 1000 characters
							chronlyfilter := fmt.Sprintf("%v", usrInput)
							charcountfilter := fmt.Sprintf("%v", strings.Join(strings.Fields(chronlyfilter), ""))
							chrcount := uniseg.GraphemeClusterCount(charcountfilter)

							if chrcount < 1000 {

								// notify galpt
								notifyCreator = true

								if isCodex {

									for fileIdx := range msgAttachment {

										// Get the file and write it to memory
										getFile, err := httpclient.Get(msgAttachment[fileIdx].URL)
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}

											break
										}

										// ==================================
										// Create a new uid.txt file
										createcdxRespFile, err := memFS.Create(fmt.Sprintf("./OpenAI/%v.txt", userID))
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}

											break
										}

										// Write to the file
										writecdxRespFile, err := io.Copy(createcdxRespFile, getFile.Body)
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}

											getFile.Body.Close()

											break
										}

										getFile.Body.Close()

										if err := createcdxRespFile.Close(); err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}

											break
										}

										winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", createcdxRespFile.Name(), (writecdxRespFile / Kilobyte), (writecdxRespFile / Megabyte))
										fmt.Println(winLogs)

										if len(universalLogs) >= universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
										}

										// check input file md5
										readcdxFile, err := afero.ReadFile(memFS, createcdxRespFile.Name())
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}

											break
										}

										usrInput = fmt.Sprintf("%v", string(readcdxFile))

									}

								}

								totalWords := strings.Fields(usrInput)
								if !isCodex {

									if 6 <= len(totalWords) {
										model = "text-davinci-002"
									} else if 3 <= len(totalWords) && len(totalWords) <= 5 {
										model = "text-curie-001"
									} else {
										model = "text-ada-001"
									}

								}

								c := gogpt.NewClient(apiKey)
								ctx := context.Background()

								// content filter check
								var (
									maxTokensFilter = 1
									tempFilter      = float32(0.0)
									topPFilter      = float32(0)
									nFilter         = 1
									logProbsFilter  = 10
									usrInputFilter  = ""
								)
								usrInputFilter = fmt.Sprintf("%v\n--\nLabel:", usrInput)

								reqfilter := gogpt.CompletionRequest{
									MaxTokens:        maxTokensFilter,
									Prompt:           usrInputFilter,
									Echo:             false,
									Temperature:      tempFilter,
									TopP:             topPFilter,
									N:                nFilter,
									LogProbs:         logProbsFilter,
									PresencePenalty:  float32(0),
									FrequencyPenalty: float32(0),
								}
								respfilter, err := c.CreateCompletion(ctx, "content-filter-alpha", reqfilter)
								if err != nil {
									fmt.Println(" [ERROR] ", err)

									if len(universalLogs) >= universalLogsLimit {
										universalLogs = nil
									} else {
										universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
									}

									s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("**Ei's Answer**\n\n[ERROR]\n%v", err), m.Reference())
									return
								}

								if respfilter.Choices[0].Text == "2" {
									s.ChannelMessageSendReply(m.ChannelID, "I've detected that the generated response could be sensitive or unsafe.\nRest assured, I won't send it back to you.", m.Reference())

									// decrease the counter to allow next request
									currReqPerMin = currReqPerMin - 1

									return
								} else if respfilter.Choices[0].Text == "1" || respfilter.Choices[0].Text == "0" {

									req := gogpt.CompletionRequest{
										MaxTokens:        allowedTokens,
										Prompt:           usrInput,
										Echo:             false,
										Temperature:      tempptr,
										TopP:             toppptr,
										N:                nvalptr,
										LogProbs:         openaiLogprobs,
										PresencePenalty:  openaiPresPen,
										FrequencyPenalty: openaiFreqPen,
										BestOf:           openaiBestOf,
									}
									resp, err := c.CreateCompletion(ctx, model, req)
									if err != nil {
										fmt.Println(" [ERROR] ", err)

										if len(universalLogs) >= universalLogsLimit {
											universalLogs = nil
										} else {
											universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
										}

										s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("**Ei's Answer**\n\n[ERROR]\n%v", err), m.Reference())
										return
									}

									respEdited = strings.ReplaceAll(resp.Choices[0].Text, "\n", " ")
									// totalRespWords := strings.Fields(resp.Choices[0].Text)
									// wordCount = len(totalRespWords)
									charOnly := fmt.Sprintf("%v", strings.Join(strings.Fields(respEdited), ""))
									charCount = uniseg.GraphemeClusterCount(charOnly)

									if 6 <= len(totalWords) {

										// cost for "davinci"
										costCount = (float64((uniseg.GraphemeClusterCount(resp.Choices[0].Text) / 4)) * (0.0600 / 1000))
									} else if 3 <= len(totalWords) && len(totalWords) <= 5 {

										// cost for "curie"
										costCount = (float64((uniseg.GraphemeClusterCount(resp.Choices[0].Text) / 4)) * (0.0060 / 1000))
									} else {

										// cost for "ada"
										costCount = (float64((uniseg.GraphemeClusterCount(resp.Choices[0].Text) / 4)) * (0.0008 / 1000))
									}

									// get time elapsed data
									execTime := time.Since(codeExec)

									// Create the embed templates.
									msginfoField := discordgo.MessageEmbedField{
										Name:   "Message Info",
										Value:  fmt.Sprintf("ID: `%v` | <#%v>", m.ID, m.ChannelID),
										Inline: false,
									}
									timeElapsedField := discordgo.MessageEmbedField{
										Name:   "Processing Time",
										Value:  fmt.Sprintf("`%v`", execTime),
										Inline: true,
									}
									costField := discordgo.MessageEmbedField{
										Name:   "Operational Cost",
										Value:  fmt.Sprintf("```\n• mode: %v\n• model: %v\n• chars: %v\n• tokens: %v\n• cost: $%.4f/1k tokens\n```", mode, resp.Model, charCount, (uniseg.GraphemeClusterCount(resp.Choices[0].Text) / 4), costCount),
										Inline: true,
									}
									messageFields := []*discordgo.MessageEmbedField{&msginfoField, &timeElapsedField, &costField}

									aoiEmbedFooter := discordgo.MessageEmbedFooter{
										Text: fmt.Sprintf("%v's Server Time • %v", botName, time.Now().UTC().Format(time.RFC850)),
									}

									aoiEmbedAuthor := discordgo.MessageEmbedAuthor{
										URL:     fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
										Name:    fmt.Sprintf("%v#%v", m.Author.Username, m.Author.Discriminator),
										IconURL: fmt.Sprintf("%v", m.Author.AvatarURL("4096")),
									}

									aoiEmbeds := discordgo.MessageEmbed{
										Title:  "Intelli-Ei",
										Color:  0x7581eb,
										Footer: &aoiEmbedFooter,
										Fields: messageFields,
										Author: &aoiEmbedAuthor,
									}

									s.ChannelMessageSendEmbed(m.ChannelID, &aoiEmbeds)

									// =========================
									// Create the embed template for notifyCreator.
									notifaskedQField := discordgo.MessageEmbedField{
										Name:   "Asked Question",
										Value:  fmt.Sprintf("%v", usrInput),
										Inline: false,
									}
									notifmessageFields := []*discordgo.MessageEmbedField{&msginfoField, &timeElapsedField, &costField, &notifaskedQField}

									notifyEmbeds := discordgo.MessageEmbed{
										Title:  "Notify Intelli-Ei",
										Color:  0x7581eb,
										Footer: &aoiEmbedFooter,
										Fields: notifmessageFields,
										Author: &aoiEmbedAuthor,
									}

									if notifyCreator {
										// send a copy to galpt
										channel, err := s.UserChannelCreate(staffID[0])
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}
											notifyCreator = false
											return
										}
										_, err = s.ChannelMessageSendEmbed(channel.ID, &notifyEmbeds)
										if err != nil {
											fmt.Println(" [ERROR] ", err)

											if len(universalLogs) >= universalLogsLimit {
												universalLogs = nil
											} else {
												universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
											}
											notifyCreator = false
											return
										}

										// notify galpt
										notifyCreator = false
									}

									if sendCodex {

										// ==================================
										// Create a new reply.txt
										createReplyFile, err := memFS.Create("./OpenAI/reply.txt")
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
											writeReplyFile, err := createReplyFile.WriteString(fmt.Sprintf("%v", resp.Choices[0].Text))
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
												if err := createReplyFile.Close(); err != nil {
													fmt.Println(" [ERROR] ", err)

													if len(universalLogs) >= universalLogsLimit {
														universalLogs = nil
													} else {
														universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
													}

													return
												} else {
													winLogs = fmt.Sprintf(" [DONE] `%v` file has been created. \n >> Size: %v KB (%v MB)", createReplyFile.Name(), (writeReplyFile / Kilobyte), (writeReplyFile / Megabyte))
													fmt.Println(winLogs)

													if len(universalLogs) >= universalLogsLimit {
														universalLogs = nil
													} else {
														universalLogs = append(universalLogs, fmt.Sprintf("\n%v", winLogs))
													}
												}
											}
										}

										readOutput, err := afero.ReadFile(memFS, fmt.Sprintf("%v", createReplyFile.Name()))
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

										s.ChannelFileSend(m.ChannelID, fmt.Sprintf("%v.%v.txt", userID, execTime), reader)

									} else {
										s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("%v", resp.Choices[0].Text), m.Reference())
									}

								} else {
									return
								}

							} else {
								s.ChannelMessageSendReply(m.ChannelID, "You are not allowed to ask more than 1000 characters.", m.Reference())
							}

							// decrease the counter to allow next request
							currReqPerMin = currReqPerMin - 1

						} else {
							// if there's a user using the AI right now,
							// wait until the request is finished.
							s.ChannelMessageSendReply(m.ChannelID, "There's a user using the AI right now.\nPlease wait until the process is finished.", m.Reference())
						}

					}

				}

			}

		}
	}

}
