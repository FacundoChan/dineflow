package config

import (
	"github.com/spf13/viper"
	"strings"
)

func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	return viper.ReadInConfig()
}
