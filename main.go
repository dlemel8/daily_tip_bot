package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use: "daily_tip_bot",
	}
	webServiceCmd = &cobra.Command{
		Use: "web_service",
		RunE: func(cmd *cobra.Command, args []string) error {
			signingSecret := viper.GetString("slack_signing_secret")
			http.HandleFunc("/slack", newSlashCommandHandler(signingSecret))

			log.Info("Server listening")
			port := viper.GetInt("port")
			return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}
	scheduledTipsCmd = &cobra.Command{
		Use: "scheduled_tips",
		RunE: func(cmd *cobra.Command, args []string) error {
			botToken := viper.GetString("slack_bot_token")
			fmt.Printf("bot token len is %d\n", len(botToken))
			return nil
		},
	}
)

func init() {
	viper.AutomaticEnv()
	rootCmd.AddCommand(webServiceCmd)
	rootCmd.AddCommand(scheduledTipsCmd)
}

func main() {
	cmd, _, _ := rootCmd.Find(os.Args[1:])
	if cmd != nil && cmd.Name() == rootCmd.Name() {
		// set default command
		args := append([]string{webServiceCmd.Name()}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
