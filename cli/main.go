package main

import (
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"

	iterate "github.com/njnygaard/iterate-vault/vault-iterator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"function": "main",
	})

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

	args := os.Args[1:]

	switch args[0] {
	case "get":
		fallthrough
	case "g":
		logger.Info("get")
		if args[1] == "" {
			logger.Error("path required")
			return
		}else{
			var root iterate.Folder
			root.Init()
			err := iterate.Find(args[1], cfg, &root, 0)
			if err != nil {
				logger.Error(err)
			}
			spew.Dump(root)
		}
	default:
		logger.Error("unrecognized command")
		return
	}
}
