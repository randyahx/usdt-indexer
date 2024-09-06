package main

import (
	"contract-indexer/app"
	config "contract-indexer/config"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func main() {
	// Viper flags
	initFlags()

	// Initialize config from path provided
	var configFilePath string
	switch {
	case viper.GetString(config.FlagConfigPath) != "":
		configFilePath = viper.GetString(config.FlagConfigPath)
	case os.Getenv(config.ENVConfigPath) != "":
		configFilePath = os.Getenv(config.ENVConfigPath)
	default:
		configFilePath = config.DefaultConfigPath
	}
	cfg, err := config.ParseConfigFromFile(configFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to parse config file when initializing app, err=%s", err))
	}

	// Initialize application
	indexerApp := app.NewApp(cfg)
	indexerApp.Start()

	select {}
}

func initFlags() {
	flag.String(config.FlagConfigPath, "", "config file path")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Errorf("failed to init flags when initializing app, err=%s", err))
	}
}

func printUsage() {
	fmt.Print("usage: ./usdt-indexer --config-path configFile\n")
}
