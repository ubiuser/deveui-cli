package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NickGowdy/deveui-cli/client"
	"github.com/NickGowdy/deveui-cli/processor"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	BaseURL               string        `envconfig:"BASE_URL" default:"http://europe-west1-machinemax-dev-d524.cloudfunctions.net"`
	MaxConcurrent         int           `envconfig:"MAX_CONCURRENT" default:"10"`
	CodeRegistrationLimit int           `envconfig:"CODE_REGISTRATION_LIMIT" default:"100"`
	Timeout               time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func main() {
	cfg := mustReadConfig()
	loraWAN := mustCreateClient(cfg.BaseURL, cfg.Timeout)
	proc := processor.New(cfg.CodeRegistrationLimit, cfg.MaxConcurrent, loraWAN)

	ctx, cancel := context.WithCancel(context.Background())

	go proc.Start(ctx, cancel)

	handleGracefulShutdown(ctx, cancel)
}

func mustReadConfig() config {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		panic("failed to process env vars: " + err.Error())
	}

	return cfg
}

func mustCreateClient(baseURL string, timeout time.Duration) *client.LoraWAN {
	loraWAN, err := client.NewLoraWAN(baseURL, timeout)
	if err != nil {
		panic("failed to create loraWAN client: " + err.Error())
	}

	return loraWAN
}

func handleGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	const (
		shutdownTimeout = 5 * time.Second

		shutdownTimeoutExpired = 1
		forcedExit             = 2
	)

	// create interrupt signal listener, use the app's context
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	// block and listen for the interrupt signal
	<-ctx.Done()

	log.Println("shutting down")

	// let the ctx chain know that the app is terminating
	cancel()

	go func() {
		// restart listening for the force exist signal, any new context will do here
		// e.g. press Ctrl+C again to force
		ctx2, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-time.After(shutdownTimeout):
			log.Println("shutdown timeout expired")

			os.Exit(shutdownTimeoutExpired)
		case <-ctx2.Done():
			log.Println("forced exit")

			os.Exit(forcedExit)
		}
	}()
}
