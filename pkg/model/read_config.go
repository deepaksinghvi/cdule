package model

import (
	"github.com/deepaksinghvi/cdule/pkg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readConfig(param []string) (*pkg.CduleConfig, error) {
	v := viper.New()
	v.AddConfigPath(param[0]) //"./resources"
	v.SetConfigName(param[1]) // "config"
	v.AutomaticEnv()
	v.SetConfigType("yml")

	var cduleConfig pkg.CduleConfig
	if err := v.ReadInConfig(); err != nil {
		log.Error("Error reading config file ", err)
		return nil, err
	}
	err := v.Unmarshal(&cduleConfig)
	if err != nil {
		log.Error("Unable to read into CduleConfig ", err)
		return nil, err
	}
	return &cduleConfig, nil
}
