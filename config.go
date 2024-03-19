package main

import (
	"github.com/spf13/viper"
	"log"
)

func initConfig() {
	viper.SetConfigName("config") // имя файла конфигурации без расширения
	viper.SetConfigType("yaml")   // тип файла конфигурации
	viper.AddConfigPath(".")      // путь к файлу конфигурации

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.AutomaticEnv()
}
