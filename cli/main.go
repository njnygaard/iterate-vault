package main

import  (
	"github.com/sirupsen/logrus"
	iterate "github.com/njnygaard/iterate-vault/vault-iterator"
)

func main(){
	logger := logrus.WithFields(logrus.Fields{
		"status":  "working",
		"handler": "handleStripeWebhook",
	})
	logger.Info("Something Cool")

	err := iterate.Hello()
	if err != nil {
		logger.Info(err)
	}
}
