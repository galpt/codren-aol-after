package main

import (
	"crypto/tls"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/servusdei2018/shards"
	"github.com/spf13/afero"
	xurls "mvdan.cc/xurls/v2"
)

const (
	Gigabyte      = 1 << 30
	Megabyte      = 1 << 20
	Kilobyte      = 1 << 10
	timeoutTr     = 24 * time.Hour
	memCacheLimit = 300 << 20 // 300 MB
)

var (
	httpPort        = ":7777"
	discordBotToken = "" // fill your discord bot token here

	tlsConf = &tls.Config{
		InsecureSkipVerify: true,
	}

	h1Tr = &http.Transport{
		DisableKeepAlives:      false,
		DisableCompression:     false,
		ForceAttemptHTTP2:      false,
		TLSClientConfig:        tlsConf,
		TLSHandshakeTimeout:    30 * time.Second,
		ResponseHeaderTimeout:  90 * time.Second,
		IdleConnTimeout:        90 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		MaxIdleConns:           1000,     // Prevents resource exhaustion
		MaxIdleConnsPerHost:    100,      // Increases performance and prevents resource exhaustion
		MaxConnsPerHost:        0,        // 0 for no limit
		MaxResponseHeaderBytes: 64 << 10, // 64k
		WriteBufferSize:        64 << 10, // 64k
		ReadBufferSize:         64 << 10, // 64k
	}

	httpclient = &http.Client{
		Timeout:   60 * time.Second,
		Transport: h1Tr,
	}

	universalLogs      []string
	universalLogsLimit = 100

	Mgr *shards.Manager

	statusInt   = 0
	statusSlice = []string{"idle", "online", "dnd"}

	// katInz YTDL feature
	katInzVidID  = ""
	xurlsRelaxed = xurls.Relaxed()
	botName      = "Ei"

	staffID = []string{"631418827841863712"}

	osFS        = afero.NewOsFs()
	memFS       = afero.NewMemMapFs()
	httpCache   = afero.NewHttpFs(osFS)
	mem         runtime.MemStats
	duration    = time.Now()
	ReqLogs     string
	RespLogs    string
	ConnReqLogs string
	totalMem    string
	HeapAlloc   string
	SysMem      string
	Frees       string
	NumGCMem    string
	timeElapsed string
	latestLog   string
	winLogs     string
	tempDirLoc  string

	lastMsgTimestamp   string
	lastMsgUsername    string
	lastMsgUserID      string
	lastMsgpfp         string
	lastMsgAccType     string
	lastMsgID          string
	lastMsgContent     string
	lastMsgTranslation string

	maidsanLastMsgChannelID string
	maidsanLastMsgID        string
	maidsanLowercaseLastMsg string
	maidsanEditedLastMsg    string
	maidsanTranslatedMsg    string
	maidsanBanUserMsg       string

	katInzBlacklist               []string
	katInzBlacklistReadable       string
	katInzCustomBlacklistReadable string

	maidsanLogs         []string
	maidsanLogsLimit    = 500
	maidsanLogsTemplate string
	timestampLogs       []string
	useridLogs          []string
	profpicLogs         []string
	acctypeLogs         []string
	msgidLogs           []string
	msgLogs             []string
	translateLogs       []string

	maidsanWatchCurrentUser  string
	maidsanWatchPreviousUser string
	maidsanEmojiInfo         []string
	replyremoveNewLines      string
	replyremoveSpaces        string
	replysplitEmojiInfo      []string
	customEmojiReply         string
	customEmojiDetected      bool
	customEmojiIdx           = 0
	customEmojiSlice         []string

	stickerList []string
	svstatLock  = false
	ytLock      = false
	xvLock      = false

	uaChrome = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36"

	// vars for openai
	maxReqPerMin   = 10
	currReqPerMin  = 0
	openaiBestOf   = 2
	openaiPresPen  = float32(1.5)
	openaiFreqPen  = float32(1.5)
	openaiLogprobs = 0
	notifyCreator  = false
	openAIAccess   = []string{
		"631418827841863712", // castella
		"323393785352552449", // nuke
		"411531606092677121", // jef kimi no udin
		"243660664441143297", // mdx ojtojtojt
		"742020307371425823", // sinsin
		"413608064730791936", // fred Thorian#2939
	}
	re = regexp.MustCompile("[0-9]+")
)
