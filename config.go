package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func Config() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/." + Name)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
