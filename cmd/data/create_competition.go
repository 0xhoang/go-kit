package data

import (
	"fmt"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/services"
	"github.com/spf13/cobra"
)

func createCompetition(apiSvc *services.HelloService) *cobra.Command {
	var cmd = &cobra.Command{
		Use: "create-competition",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Example: fmt.Sprintf(
			"./cli create-competition",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Start Create Competition")
			err := apiSvc.CreateCompetition()
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
