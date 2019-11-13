package main

import  "github.com/sirupsen/logrus"

func main(){
	logger := logrus.WithFields(logrus.Fields{
		"status":  "working",
		"handler": "handleStripeWebhook",
	})
	logger.Info("Something Cool")
}
