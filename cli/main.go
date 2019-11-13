package main

import (
	"io/ioutil"
	"os"

	iterate "github.com/njnygaard/iterate-vault/vault-iterator"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"status":  "working",
		"handler": "handleStripeWebhook",
	})
	logger.Info("Something Cool")

	configPath := os.Getenv("ITERATOR_CONFIG_FILE")
	if configPath == "" {
		configPath = ".auth.yaml"
	}

	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	var cfg iterate.AuthConfig

	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		logger.Error(err)
		return
	}

	err = iterate.Find("something", "something_else", cfg)
	if err != nil {
		logger.Info(err)
	}
}
