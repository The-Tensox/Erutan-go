package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	server "github.com/user/erutan_two/server"
	utils "github.com/user/erutan_two/utils"
	"golang.org/x/net/context"
)

func init() {
	log.Printf("Init package")
}

func flags() {
	flag.BoolVar(&utils.Config.DebugMode, "v", false, "enable debug logging")
	flag.StringVar(&utils.Config.Host, "h", "0.0.0.0:50051", "the chat server's host")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
}

// RunMain starts the server and everything for prod, used in tests
func RunMain() {
	ctx := utils.SignalContext(context.Background())
	var err error

	err = server.NewServer(utils.Config.Host).Run(ctx)

	if err != nil {
		utils.MessageLog(time.Now(), "<<Process>>", err.Error())
		os.Exit(1)
	}
}

func main() {
	flags()
	utils.InitializeConfig(35)
	RunMain()
}
