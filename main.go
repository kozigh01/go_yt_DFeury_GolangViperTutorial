package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

var (
	defaults = map[string]interface{}{
		"psername": "admin",
		"password": "password1",
		"host":     "localhost",
		"port":     3306,
		"database": "test",
	}
	configName  = "config.yml"
	configPaths = []string{
		"$HOME/.appname",
		".",
	}
)

func initViperDefaults() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

func main() {
	initViperDefaults()

	viper.SetConfigFile(configName)
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %q\n", err)
	}

	fmt.Printf("Username from viper: %q\n", viper.GetString("username"))
	fmt.Printf("Password from viper: %q\n", viper.GetString("password"))
	fmt.Printf("Host from viper: %q\n", viper.GetString("host"))
	fmt.Printf("Port from viper: %v\n", viper.GetInt("port"))
	fmt.Printf("Database from viper: %q\n", viper.GetString("database"))

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("could not decode config into struct:\n%v", err)
	}
	fmt.Println(config)

	// changed := false
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	err = viper.Unmarshal(&config)
	// 	if err != nil {
	// 		log.Printf("could not decode config into struct after changed:\n%v", err)
	// 	}
	// 	changed = true
	// })
	// for !changed {
	// 	time.Sleep(time.Second)
	// 	fmt.Printf("Config struct: %v\n", config)
	// }

	changed := make(chan Config)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Printf("could not decode config into struct after changed:\n%v", err)
		}
		changed <- config
	})
	fmt.Printf("Config struct: %v\n", <-changed)
}
