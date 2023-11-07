package main

import (
	"fmt"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func main() {
	cfg := config.ReadConfigAndArg()
	logger, _, err := must.NewLogger(cfg.SentryDSN, cfg.ServiceName+"-local")
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	helloSvc := services.NewGokitService(
		logger,
		cfg,
		nil,
		nil,
	)

	rootCmd := newRootCmd(helloSvc)

	if err := rootCmd.Execute(); err != nil {
		fmt.Print(common.Red)
		fmt.Println("Execute", err.Error())
		fmt.Print(common.Reset)
		os.Exit(1)
	}
}

func newRootCmd(apiSvc *services.GokitService) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "./cli",
	}

	rootCmd.AddCommand(
		createCompetition(apiSvc),
		createTeam(apiSvc),
	)

	return rootCmd
}
