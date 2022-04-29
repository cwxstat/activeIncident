/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cwxstat/activeIncident/constants"
	"github.com/cwxstat/activeIncident/db"
	"github.com/cwxstat/activeIncident/dbpop"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "activeIncident",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		for {

			func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				as, err := db.NewActiveIncidentServer(ctx)
				if err != nil {
					log.Println(err)
					time.Sleep(constants.ErrorBackoff)
					return
				}

				a, err := dbpop.NewActiveIncidentEntry()
				if err != nil {
					log.Println(err)
				}

				err = as.AddEntry(ctx, a)
				if err != nil {
					log.Println(err)
					return
				}
				if err := as.Disconnect(ctx); err != nil {
					log.Println("as.Disconnect: ", err)
				}
				log.Println("entry added")

			}()

			time.Sleep(constants.RefreshRate)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.activeIncident.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
