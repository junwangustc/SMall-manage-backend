package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/junwangustc/ustclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *log.Logger
var logfile string
var configPath string
var GitCommit string
var BuildTime string
var onlyShowVersion bool

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.BoolVar(&onlyShowVersion, "version", false, "just show version")
	flag.StringVar(&logfile, "logfile", "/tmp/SMall-manage-backend.log", "output log file")
	flag.StringVar(&configPath, "config", "/tmp/SMall-manage-backend.toml", "Small-manage-backend config file path ")
	flag.Parse()
}
func initConfig() {
	output := &lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
	}
	log.SetOutput(output)
	globalLogger = log.New(output, "", log.Ldefault)
}
func main() {
	if GitCommit == "" {
		GitCommit = "UNKNOW"
	}
	if BuildTime == "" {
		BuildTime = "UNKNOW"
	}
	fmt.Printf("GitCommit: %s, BuildTime: %s\n", GitCommit, BuildTime)
	if onlyShowVersion {
		return
	}
	initConfig()
	conf, err := ParseConfig(configPath)
	if err != nil {
		log.Println("[Error]configfile error", err)
		return
	}
	log.Println("parse config success,the parse config value is ", conf)
	srv, err := NewServer(conf, globalLogger)
	if err != nil {
		log.Println("[Error]new server", err)
		return
	}
	if err := srv.Open(); err != nil {
		log.Println("[Error]open server fail", err)
		return
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Println("Listening for signals")
	select {
	case <-signalCh:
		log.Println("Signal received, shutdown...")
	}
}
