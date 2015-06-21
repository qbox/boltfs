package main

import (
	"net/http"
	"runtime"

	"qbox.us/cc/config"

	"github.com/qiniu/http/restrpc.v1"
	"github.com/qiniu/log.v1"

	"qiniu.com/qfusegate.v1"
)

// ---------------------------------------------------------------------------

type Config struct {
	Gate qfusegate.Config `json:"gate"`

	BindHost   string `json:"bind_host"`
	MaxProcs   int    `json:"max_procs"`
	DebugLevel int    `json:"debug_level"`
}

func main() {

	// Load Config

	config.Init("f", "qiniu", "qfusegate.conf")

	var conf Config
	if err := config.Load(&conf); err != nil {
		log.Fatal("config.Load failed:", err)
	}
	log.Info("config:", conf)

	// General Settings

	runtime.GOMAXPROCS(conf.MaxProcs)
	log.SetOutputLevel(conf.DebugLevel)

	// new Service

	service, err := qfusegate.New(&conf.Gate)
	if err != nil {
		log.Fatal("qfusegate.New failed:", err)
	}

	// run Service

	router := restrpc.Router{
		PatternPrefix: "v1",
	}
	log.Info("Starting qfusegate ...")
	err = http.ListenAndServe(conf.BindHost, router.Register(service))
	log.Fatal("http.ListenAndServe(qfusegate):", err)
}

// ---------------------------------------------------------------------------

