package agent

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var globalConfig map[string]string

func init() {
	globalConfig = make(map[string]string)
}

func readConfig() error {
	log.Printf("[agent] reading configuration from: %s\n", configFile)
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		// TODO
		return err
	}
	return parseConfig()
}

// TODO move this to pkg directory as config package
func parseConfig() error {
	hostname := viper.GetString("hostname")
	if hostname == "" {
		return fmt.Errorf("[config] please configure 'hostname' option in config file %s", configFile)
	}
	globalConfig["hostname"] = viper.GetString("hostname")

	statsd_port := viper.GetString("statsd_port")
	if statsd_port == "" {
		return fmt.Errorf("[config] please configure 'statsd_port' option in config file %s", configFile)
	}
	globalConfig["statsd_port"] = viper.GetString("statsd_port")

	tornimo_put := viper.GetString("tornimo_put")
	if tornimo_put == "" {
		return fmt.Errorf("[config] please configure 'tornimo_put' option in config file %s", configFile)
	}
	globalConfig["tornimo_put"] = viper.GetString("tornimo_put")

	tornimo_token := viper.GetString("tornimo_token")
	if tornimo_token == "" {
		return fmt.Errorf("[config] please configure 'tornimo_token' option in config file %s", configFile)
	}
	globalConfig["tornimo_token"] = viper.GetString("tornimo_token")

	for k, v := range globalConfig {
		log.Printf("[config] %s = %s\n", k, v)
	}
	return nil
}
