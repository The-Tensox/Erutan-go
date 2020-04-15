package main

import (
	"fmt"
	"github.com/The-Tensox/erutan/cfg"
	"github.com/The-Tensox/erutan/server"
	"github.com/The-Tensox/erutan/utils"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
)

func init() {
	log.Printf("Init package")
}


// RunMain starts the server and everything for prod, used in tests
func RunMain() {
	ctx := utils.SignalContext(context.Background())
	var err error

	err = server.NewServer(fmt.Sprintf("%s:%s", cfg.Global.Server.Host, cfg.Global.Server.Port)).Run(ctx)

	if err != nil {
		utils.MessageLog(time.Now(), "<<Process>>", err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	cfg.Global = cfg.Get()
	utils.ServerLogf(time.Now(), "Starting server with config %v", cfg.Global)
	RunMain()
}

