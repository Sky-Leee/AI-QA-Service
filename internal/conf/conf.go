package conf

import "github.com/spf13/viper"

func InitConf(confPath string) {
	viper.SetConfigFile(confPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
}
