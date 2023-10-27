package main

import (
	"fmt"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/spf13/cobra"
)

func createTeam(apiSvc *services.GokitService) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "create-team",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Example: fmt.Sprintf(
			"./cli create-team",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Start Create Team")
			err := apiSvc.CreateTeam()
			if err != nil {
				fmt.Print(common.Red)
				//spew.Dump(rs)
				fmt.Printf("Error => %s\n", err.Error())
				fmt.Print(common.Reset)
				return nil
			}

			fmt.Print(common.Blue)
			//spew.Dump(rs)
			fmt.Printf("Create Succesfully\n")
			fmt.Print(common.Reset)
			return nil
		},
	}

	return cmd
}
