package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	server "github.com/The-Tensox/erutan/server"
	utils "github.com/The-Tensox/erutan/utils"
	"golang.org/x/net/context"
)

func init() {
	log.Printf("Init package")
}

func flags() {
	flag.BoolVar(&utils.Config.DebugMode, "d", false, "enable debug logging")
	flag.StringVar(&utils.Config.Host, "h", "0.0.0.0:50051", "the server's host")
	flag.Float64Var(&utils.Config.GroundSize, "g", 20, "ground size")
	flag.Float64Var(&utils.Config.TimeScale, "t", 1, "the server's time scale")
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
	RunMain()
}
