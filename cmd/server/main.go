package main

import (
	"fmt"
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/server"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
)



// RunMain starts the server and everything for prod, used in tests
func RunMain() {
	log.Init()
	defer log.Zap.Sync()
	log.Zap.Info( "Starting server with config", zap.Any("config", cfg.Get()))
	ctx := signalContext(context.Background())
	var err error

	err = server.NewServer(fmt.Sprintf("%s:%s", cfg.Get().Server.Host, cfg.Get().Server.Port)).Run(ctx)

	if err != nil {
		log.Zap.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	RunMain()
}

func signalContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Zap.Info("listening for shutdown signal")
		<-sigs
		log.Zap.Info("shutdown signal received")
		signal.Stop(sigs)
		close(sigs)
		cancel()
	}()

	return ctx
}