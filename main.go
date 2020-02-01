package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
)

var (
	debugMode bool
	host      string
	password  string
)

func init() {
	log.Printf("Init package")
}

func flags() {
	flag.BoolVar(&debugMode, "v", false, "enable debug logging")
	flag.StringVar(&host, "h", "0.0.0.0:50051", "the chat server's host")
	flag.StringVar(&password, "p", "", "the chat server's password")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
}

// RunMain starts the server and everything for prod, used in tests
func RunMain() {
	ctx := SignalContext(context.Background())
	var err error

	err = NewServer(host, password).Run(ctx)

	if err != nil {
		MessageLog(time.Now(), "<<Process>>", err.Error())
		os.Exit(1)
	}
}

func main() {
	flags()
	RunMain()
}
