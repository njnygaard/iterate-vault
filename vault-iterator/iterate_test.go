package iterate

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func setup() (c AuthConfig, err error) {
	configPath := os.Getenv("ITERATOR_CONFIG_FILE")
	if configPath == "" {
		configPath = ".auth.yaml"
	}

	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(bs, &c); err != nil {
		return
	}

	return
}

func TestFind_Leaf(t *testing.T) {

	logger := logrus.WithFields(logrus.Fields{
		"test": "TestFind_Leaf",
	})

	cfg, err := setup()
	if err != nil {
		t.Errorf("configuration failed with error = %q", err)
	}

	var root Leaf
	root.Init()

	if err := Find("secret/deployments/k8s/default/services/cloudmanager/config/api", cfg, &root, 0); err != nil {
		t.Errorf("Get() errored with err = %q", err)
	}

	// https://stackoverflow.com/a/28384502/1236359
	// Dereference the pointer to the map first, then index it.
	if (*root.data)["UNINSTALL_SCRIPT"] != "uninstall_1.0.1.sh" {
		t.Error("unexpected value error")
	}

	logger.Info("find leaf passed")
}

func TestFind_Folder(t *testing.T) {

	logger := logrus.WithFields(logrus.Fields{
		"test": "TestFind_Folder",
	})

	cfg, err := setup()
	if err != nil {
		t.Errorf("configuration failed with error = %q", err)
	}

	var root Folder
	root.Init()

	if err := Find("secret/deployments/k8s/default/services/cloudmanager", cfg, &root, 0); err != nil {
		t.Errorf("Get() errored with err = %q", err)
	}

	if (*(*(*root.childFolders)[0].childLeaves)[0].data)["UNINSTALL_SCRIPT"] != "uninstall_1.0.1.sh" {
		t.Error("unexpected value error")
	}

	logger.Info("find folder passed")
}

func TestMove_Folder(t *testing.T) {

	logger := logrus.WithFields(logrus.Fields{
		"test": "TestMove_Folder",
	})

	source := "secret/deployments/k8s/default/services/cloudmanager/"
	destination := "secret/testing/k8s/default/services/cloudmanager/"

	cfg, err := setup()
	if err != nil {
		t.Errorf("configuration failed with error = %q", err)
	}

	var root Folder
	root.Init()
	var populated Folder
	populated.Init()

	if err := Move(source, destination, cfg, &root, 0); err != nil {
		t.Errorf("Move() errored with err = %q", err)
	}

	if err := Find(destination, cfg, &populated, 0); err != nil {
		t.Errorf("Get() errored with err = %q", err)
	}

	if (*(*(*populated.childFolders)[0].childLeaves)[0].data)["UNINSTALL_SCRIPT"] != "uninstall_1.0.1.sh" {
		t.Error("unexpected value error")
	}

	logger.Info("move leaf passed")
}

//func TestMove_Leaf(t *testing.T) {
//
//	logger := logrus.WithFields(logrus.Fields{
//		"test": "TestMove_Folder",
//	})
//
//	cfg, err := setup()
//	if err != nil {
//		t.Errorf("configuration failed with error = %q", err)
//	}
//
//	var root Folder
//	root.Init()
//
//	if err := Find("secret/deployments/k8s/default/services/cloudmanager", cfg, &root, 0); err != nil {
//		t.Errorf("Get() errored with err = %q", err)
//	}
//
//	if (*(*(*root.childFolders)[0].childLeaves)[0].data)["UNINSTALL_SCRIPT"] != "uninstall_1.0.1.sh" {
//		t.Error("unexpected value error")
//	}
//
//	logger.Info("move folder passed")
//}
