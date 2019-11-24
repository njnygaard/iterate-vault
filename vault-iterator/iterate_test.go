package iterate

import (
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func TestFind(t *testing.T) {

	logger := logrus.WithFields(logrus.Fields{
		"status":  "working",
		"handler": "handleStripeWebhook",
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

	var cfg AuthConfig

	if err := yaml.Unmarshal(bs, &cfg); err != nil {
		logger.Error(err)
		return
	}

	//var want error

	/*
		vault kv list -format=json secret/deployments/k8s/default/services/cloudmanager/config/
		[
		  "api",
		  "kubernetes/",
		  "ui"
		]

		vault kv get -format=json secret/deployments/k8s/default/services/cloudmanager/config/
		No value found at secret/data/deployments/k8s/default/services/cloudmanager/config

		vault kv list -format=json secret/deployments/k8s/default/services/cloudmanager/config/kubernetes
		[
		  "cr",
		  "crd"
		]

		List Path
		secret/metadata/deployments/k8s/minikube/services/cloudmanager/config/kubernetes/cr

		Get Path
		secret/data/deployments/k8s/minikube/services/cloudmanager/config/kubernetes/cr
	*/
	//          secret/metadata/deployments/k8s/minikube/services/cloudmanager/config/kubernetes/cr
	var root Folder
	root.init()

	if err := Find("secret/deployments/k8s/default/services/cloudmanager", cfg, &root, 0); err != nil {
		t.Errorf("Get() errored with err = %q", err)
	} else {
		spew.Dump(root)
	}
}
