package cmd

import (
	"MyStonks-go/internal/server"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the server",
	Long:  "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("config: %v", viper.AllSettings())
		server.StartServer(viper.GetViper())
	},
}

func init() {
	{
		startCmd.PersistentFlags().String("server.RunMode", "debug", "server run mode")
		if err := viper.BindPFlag("server.RunMode", startCmd.PersistentFlags().Lookup("server.RunMode")); err != nil {
			log.Panic().Err(err).Send()
		}
		fmt.Println("xxxxxxx", viper.GetString("server.RunMode"))
		startCmd.PersistentFlags().Int("server.HttpPort", 8000, "server http port")
		if err := viper.BindPFlag("server.HttpPort", startCmd.PersistentFlags().Lookup("server.HttpPort")); err != nil {
			log.Panic().Err(err).Send()
		}
		startCmd.PersistentFlags().Int("server.ReadTimeout", 60, "server read timeout")
		if err := viper.BindPFlag("server.ReadTimeout", startCmd.PersistentFlags().Lookup("server.ReadTimeout")); err != nil {
			log.Panic().Err(err).Send()
		}
		startCmd.PersistentFlags().Int("server.WriteTimeout", 60, "server write timeout")
		if err := viper.BindPFlag("server.WriteTimeout", startCmd.PersistentFlags().Lookup("server.WriteTimeout")); err != nil {
			log.Panic().Err(err).Send()
		}
	}
	{
		startCmd.PersistentFlags().String("mysql.host", "localhost", "mysql host")
		if err := viper.BindPFlag("mysql.host", startCmd.PersistentFlags().Lookup("mysql.host")); err != nil {
			log.Panic().Err(err).Send()
		}

		startCmd.PersistentFlags().Int("mysql.port", 3306, "mysql port")
		if err := viper.BindPFlag("mysql.port", startCmd.PersistentFlags().Lookup("mysql.port")); err != nil {
			log.Panic().Err(err).Send()
		}

		startCmd.PersistentFlags().String("mysql.user", "root", "mysql user")
		if err := viper.BindPFlag("mysql.user", startCmd.PersistentFlags().Lookup("mysql.user")); err != nil {
			log.Panic().Err(err).Send()
		}
		startCmd.PersistentFlags().String("mysql.password", "root", "mysql password")
		if err := viper.BindPFlag("mysql.password", startCmd.PersistentFlags().Lookup("mysql.password")); err != nil {
			log.Panic().Err(err).Send()
		}

		startCmd.PersistentFlags().String("mysql.database", "mystonks", "mysql database")
		if err := viper.BindPFlag("mysql.database", startCmd.PersistentFlags().Lookup("mysql.database")); err != nil {
			log.Panic().Err(err).Send()
		}
	}
	rootCmd.AddCommand(startCmd)
}
