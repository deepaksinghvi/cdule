package model

import (
	"github.com/deepaksinghvi/cdule/pkg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readConfig(param []string) (*pkg.CduleConfig, error) {
	viper.AddConfigPath(param[0]) //"./resources"
	viper.SetConfigName(param[1]) // "config"
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var cduleConfig pkg.CduleConfig
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file ", err)
		return nil, err
	}
	err := viper.Unmarshal(&cduleConfig)
	if err != nil {
		log.Error("Unable to read into CduleConfig ", err)
		return nil, err
	}
	return &cduleConfig, nil
}
