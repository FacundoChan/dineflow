package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return err
}

func newViperConfig() error {
	relativePath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relativePath)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (relativePath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	_, here, _, _ := runtime.Caller(0)
	relativePath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	// fmt.Printf("caller from: %s, here: %s, relativePath: %s\n", callerPwd, here, relativePath)
	return relativePath, err
}
