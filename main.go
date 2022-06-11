package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	totalmem "github.com/pbnjay/memory"
	"github.com/servusdei2018/shards"
)

// =========================================
// The main function of Katheryne bot
func main() {

	// Automatically set GOMAXPROCS to the number of your CPU cores.
	// Increase performance by allowing Golang to use multiple processors.
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs) // Sets the GOMAXPROCS value
	totalMem = fmt.Sprintf("Available OS Memory: %v MB | %v GB", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
	fmt.Println()
	totalcpu := fmt.Sprintf("Available CPUs: %v", numCPUs)

	fmt.Println(totalcpu)
	fmt.Println(totalMem)
	fmt.Println(lastMsgTimestamp)
	fmt.Println(lastMsgUsername)
	fmt.Println(lastMsgUserID)
	fmt.Println(lastMsgpfp)
	fmt.Println(lastMsgAccType)
	fmt.Println(lastMsgID)
	fmt.Println(lastMsgContent)
	fmt.Println(lastMsgTranslation)
	fmt.Println(katInzBlacklistReadable)
	fmt.Println(katInzCustomBlacklistReadable)

	// run http server
	go proxyServer()
	fmt.Println("HTTP server runs on port " + httpPort)

	// Create the logs folder
	osFS.RemoveAll("./logs/")
	createLogFolder := osFS.MkdirAll("./logs/", 0777)
	if createLogFolder != nil {
		fmt.Println(" [ERROR] ", createLogFolder)
	}
	fmt.Println(` [DONE] New "logs" folder has been created. \n >> `, createLogFolder)

	createDBFolder := osFS.MkdirAll("./db/", 0777)
	if createDBFolder != nil {
		fmt.Println(" [ERROR] ", createDBFolder)
	}
	fmt.Println(` [DONE] New "db" folder has been created. \n >> `, createDBFolder)

	// Create the ./cache/ folder
	osFS.RemoveAll("./cache/")
	createCacheFolder := osFS.MkdirAll("./cache/", 0777)
	if createCacheFolder != nil {
		fmt.Println(" [ERROR] ", createCacheFolder)
	}
	fmt.Println(` [DONE] New "cache" folder has been created. \n >> `, createCacheFolder)

	// Get the latest sticker list
	fmt.Println(" Fetching sticker list. Please wait...")
	getStickers, err := httpclient.Get("https://2.castella.network/stickers.txt")
	if err != nil {
		fmt.Println(" [getStickers] ", err)

		if len(universalLogs) >= universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	bodyStickers, err := ioutil.ReadAll(bufio.NewReader(getStickers.Body))
	if err != nil {
		fmt.Println(" [bodyStickers] ", err)

		if len(universalLogs) >= universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	newstickerlist := strings.Split(string(bodyStickers), "\n")
	stickerList = append(stickerList, newstickerlist...)
	newstickerlist = nil
	fmt.Println(" Successfully fetched the sticker list.")

	// Katheryne Inazuma goroutines
	// Get the latest blocklist for dnscrypt
	getBlocklist, err := httpclient.Get("https://raw.githubusercontent.com/notracking/hosts-blocklists/master/dnscrypt-proxy/dnscrypt-proxy.blacklist.txt")
	if err != nil {
		fmt.Println(" [ERROR] ", err)

		if len(universalLogs) >= universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	bodyBlocklist, err := ioutil.ReadAll(bufio.NewReader(getBlocklist.Body))
	if err != nil {
		fmt.Println(" [ERROR] ", err)

		if len(universalLogs) >= universalLogsLimit {
			universalLogs = nil
		} else {
			universalLogs = append(universalLogs, fmt.Sprintf("\n%v", err))
		}

		return
	}

	katInzBlacklistReadable = fmt.Sprintf("\n%v\n", string(bodyBlocklist))
	katInzBlacklist = strings.Split(string(bodyBlocklist), "\n")

	// Create a new shard manager using the provided bot token.
	Mgr, err := shards.New("Bot " + discordBotToken)
	if err != nil {
		fmt.Println("[ERROR] Error creating manager,", err)
		return
	}

	// Set custom HTTP client
	Mgr.Gateway.Client = httpclient
	Mgr.Gateway.Compress = true

	// Register the messageCreate func as a callback for MessageCreate events
	Mgr.AddHandler(maidsanEmojiReact)
	Mgr.AddHandler(maidsanAutoCheck)
	Mgr.AddHandler(emojiReactions)
	Mgr.AddHandler(getUserInfo)
	Mgr.AddHandler(getCovidData)
	Mgr.AddHandler(getServerStatus)
	Mgr.AddHandler(ucoverModsDelMsg)
	Mgr.AddHandler(katInzYTDL)
	Mgr.AddHandler(katMonShowLastSender)
	Mgr.AddHandler(openAI)
	Mgr.AddHandler(xvid)

	// Register the onConnect func as a callback for Connect events.
	Mgr.AddHandler(onConnect)

	// In this example, we only care about receiving message events.
	//Mgr.RegisterIntent(discordgo.IntentsAll)
	Mgr.RegisterIntent(discordgo.IntentsAll)

	// Set the number of shards
	Mgr.SetShardCount(numCPUs)

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = Mgr.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Manager.
	fmt.Println("[INFO] Stopping shard manager...")
	Mgr.Shutdown()
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")

}
