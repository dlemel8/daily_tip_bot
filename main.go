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
		Use:   "web_server",
		Short: "listen and serve Slack slash commands. among other things, store scheduled tips",
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduledTipsStorage, err := newScheduledTipsStorage(viper.GetString("database_url"))
			if err != nil {
				return err
			}
			defer scheduledTipsStorage.close()

			signingSecret := viper.GetString("slack_signing_secret")
			http.HandleFunc("/slack", newSlashCommandHandler(signingSecret, scheduledTipsStorage))

			log.Info("Server listening")
			port := viper.GetInt("port")
			return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}
	scheduledTipsCmd = &cobra.Command{
		Use:   "scheduled_tips_sender",
		Short: "load stored scheduled tips and send them if needed",
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduledTipsStorage, err := newScheduledTipsStorage(viper.GetString("database_url"))
			if err != nil {
				return err
			}
			defer scheduledTipsStorage.close()

			botToken := viper.GetString("slack_bot_token")
			return sendScheduledTips(botToken, scheduledTipsStorage)
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
