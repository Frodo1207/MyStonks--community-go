package cmd

import (
	"MyStonks-go/internal/common/utils"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile string
var logLevel string

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `server`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("server command")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initAll)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "common.log_level", "info", "log level (trace, debug, info, warn, error, fatal, panic)")
	if err := viper.BindPFlag("common.log_level", rootCmd.PersistentFlags().Lookup("common.log_level")); err != nil {
		log.Panic().Err(err).Send()
	}
	viper.SetDefault("common.log_level", "info")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}

func initLogger() {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	dir := viper.GetString("common.deploy_dir")
	err = os.MkdirAll(dir+"/log", 0755)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	logFile := &lumberjack.Logger{
		Filename: dir + "/log/server.log",
		MaxSize:  10,
		MaxAge:   30,
	}
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true, FormatCaller: utils.CallerFormater}, zerolog.ConsoleWriter{Out: logFile, TimeFormat: time.RFC3339, NoColor: true, FormatCaller: utils.CallerFormater})
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger().Level(level)
}

func initAll() {
	initConfig()
	initLogger()
}
