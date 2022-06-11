package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	totalmem "github.com/pbnjay/memory"
)

// =========================================
// HTTP server with customizable port (the default is 7777)
func proxyServer() {

	duration := time.Now()

	// Use Gin as the HTTP router
	gin.SetMode(gin.ReleaseMode)
	ginroute := gin.Default()

	// Custom NotFound handler
	ginroute.NoRoute(func(c *gin.Context) {
		c.File("./404.html")
	})

	// print universalLogs slice
	ginroute.GET("/logs", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Total Logs: %v of %v \n • Time Elapsed: %v \n •===========================• \n • [UNIVERSAL LOGS] \n •===========================• \n \n%v \n\n", time.Now().Format(time.RFC850), totalMem, NumGCMem, len(universalLogs), universalLogsLimit, timeElapsed, universalLogs)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// Print homepage.
	ginroute.GET("/", func(c *gin.Context) {

		// Construct data-to-print in realtime
		realtimeData := fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Total Logs: %v of %v \n • Time Elapsed: %v \n •===========================• \n", time.Now().Format(time.RFC850), totalMem, NumGCMem, len(universalLogs), universalLogsLimit, timeElapsed)

		c.String(http.StatusOK, realtimeData)
	})

	// get available emoji info
	ginroute.GET("/emoji", func(c *gin.Context) {

		runtime.ReadMemStats(&mem)
		totalMem = fmt.Sprintf("%v MB (%v GB)", (totalmem.TotalMemory() / Megabyte), (totalmem.TotalMemory() / Gigabyte))
		NumGCMem = fmt.Sprintf("%v", mem.NumGC)
		timeElapsed = fmt.Sprintf("%v", time.Since(duration))
		latestLog = fmt.Sprintf("\n •===========================• \n • [SERVER STATUS] \n • Last Modified: %v \n • Total OS Memory: %v \n • Completed GC Cycles: %v \n • Time Elapsed: %v \n •===========================• \n • [AVAILABLE EMOJI LIST] \n • Total Available Emoji: %v \n •===========================• \n \n[Name —— Emoji ID —— Animated (true/false) —— Guild Name —— Guild ID]\n\n%v \n\n", time.Now().UTC().Format(time.RFC850), totalMem, NumGCMem, timeElapsed, len(maidsanEmojiInfo), maidsanEmojiInfo)

		c.String(http.StatusOK, fmt.Sprintf("%v", latestLog))

	})

	// Control Windows OS through proxy
	ginroute.StaticFS("/temp", http.Dir(os.TempDir()))
	ginroute.GET("/gettemp", func(c *gin.Context) {

		// Get the location of the TEMP dir
		tempDirLoc = fmt.Sprintf(" [DONE] Detected TEMP folder location \n >> %v", os.TempDir())
		c.String(http.StatusOK, tempDirLoc)
	})
	ginroute.GET("/deltemp", func(c *gin.Context) {

		// Delete the entire TEMP folder.
		// If it gets deleted properly, create a new TEMP folder.
		delTemp := osFS.RemoveAll(os.TempDir())
		if delTemp == nil {
			mkTemp := osFS.MkdirAll(os.TempDir(), 0777)
			if mkTemp != nil {
				winLogs = "\n • [ERROR] Failed to recreate TEMP folder. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now().Format(time.RFC850)) + "\n • Reason >> " + fmt.Sprintf("%v", mkTemp)
				c.String(http.StatusOK, winLogs)
			}
			winLogs = "\n • [DONE] TEMP folder has been cleaned. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now().Format(time.RFC850)) + "\n • Reason >> " + fmt.Sprintf("%v", mkTemp)
			c.String(http.StatusOK, winLogs)
		} else {
			winLogs = "\n • [ERROR] Failed to delete some files. \n • Timestamp >> " + fmt.Sprintf("%v", time.Now().Format(time.RFC850)) + "\n • Reason >> " + fmt.Sprintf("%v", delTemp)
			c.String(http.StatusOK, winLogs)
		}
	})

	// get data from memory
	osFS.RemoveAll("./cache/")
	osFS.MkdirAll("./cache/", 0777)
	ginroute.StaticFS("/memory", httpCache.Dir("./cache/"))

	// shared data from disk
	ginroute.Static("/yt", "./ytdl")
	ginroute.Static("/xv", "./xvids")

	// HTTP proxy server
	httpserver := &http.Server{
		Addr:              httpPort,
		Handler:           ginroute,
		TLSConfig:         tlsConf,
		MaxHeaderBytes:    64 << 10, // 64k
		ReadTimeout:       timeoutTr,
		ReadHeaderTimeout: timeoutTr,
		WriteTimeout:      timeoutTr,
		IdleTimeout:       timeoutTr,
	}
	httpserver.SetKeepAlivesEnabled(true)
	httpserver.ListenAndServe()
}
