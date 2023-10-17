package data

import (
	"github.com/0xhoang/go-kit/services"
	"github.com/spf13/cobra"
)

func NewRootCmd(apiSvc *services.HelloService) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "./cli",
	}

	rootCmd.AddCommand(
		createCompetition(apiSvc),
		createTeam(apiSvc),
	)

	return rootCmd
}
