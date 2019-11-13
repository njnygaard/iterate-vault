package iterate

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// AuthConfig is the structure of the configuration
type AuthConfig struct {
	Token     string `yaml:"token"`
	VaultAddr string `yaml:"vault_addr"`
}

// Find does stuff
func Find(key string, prop string, config AuthConfig) (err error){

	if key == "" {
		return errors.New("must provide key")
	}

	c, err := api.NewClient(&api.Config{
		Address: config.VaultAddr,
	})

	if err != nil {
		return err
	}

	c.SetToken(config.Token)

	sec, err := c.Logical().Read(key)
	if err != nil {
		return err
	}

	if sec == nil || sec.Data == nil {
		return errors.New("no data for key")
	}

	if prop == "" {
		fmt.Println("Secret data:")
		for k, v := range sec.Data {
			fmt.Printf(" - %s -> %v\n", k, v)
		}
	} else {
		fmt.Printf("%s:%s -> %v\n", key, prop, sec.Data[prop])
	}

	return
}

