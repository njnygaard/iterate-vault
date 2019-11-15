package iterate

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

// AuthConfig is the structure of the configuration
type AuthConfig struct {
	Token     string `yaml:"token"`
	VaultAddr string `yaml:"vault_addr"`
}

type Node interface {
	setName(name string)
	setData(data interface{}) (err error)
	getChildren() (children *[]Node)
	addChild(node Node)
}

type Folder struct {
	name         string
	childFolders *[]Folder
	childLeaves  *[]Leaf
}

func (f *Folder) init() {
	var folders = make([]Folder, 0)
	var leaves = make([]Leaf, 0)
	f.childFolders = &folders
	f.childLeaves = &leaves
}

func (f *Folder) setName(name string) {
	f.name = name
}

func (f *Folder) getChildren() (children *[]Node) {
	return
}

func (f *Folder) setData(data interface{}) (err error) {
	return errors.New("not implemented for folders")
}

func (f *Folder) addChild(node Node) {
	logger := logrus.WithFields(logrus.Fields{
		"status":   "working",
		"function": "Folder::addChild",
	})

	switch v := node.(type) {
	case *Folder:
		logger.Info("Folder Case")
		spew.Dump(v)
		*f.childFolders = append(*f.childFolders, *v)
	case *Leaf:
		logger.Info("Leaf Case")
		spew.Dump(v)
		*f.childLeaves = append(*f.childLeaves, *v)
	default:
		return
	}
}

type Leaf struct {
	data interface{}
	name string
}

func (l *Leaf) setName(name string) {
	l.name = name
}

func (l *Leaf) getChildren() (children *[]Node) {
	return
}

func (l *Leaf) setData(data interface{}) (err error) {
	l.data = data
	return nil
}

func (l *Leaf) addChild(node Node) {
	return
}

// Find does stuff
func Find(key string, config AuthConfig, node Node, stack int) (err error) {

	logger := logrus.WithFields(logrus.Fields{
		"status":   "working",
		"function": "Find",
	})

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

	/*
		1. list a path
		   1. if nothing, get a path
		      1. if nothing, nothing is there
		      2. if stuff, display
		   2. if stuff, look at keys
		      1. if key has `/`, list
		      2. if key does not have, get
	*/

	/*
		A friendly vault path is something like:
		vault kv list -format=json \
			secret/deployments/k8s/default/services/cloudmanager/config/kubernetes

		The List api wants that path that would get the same thing to be:
			secret/metadata/deployments/k8s/minikube/services/cloudmanager/config/kubernetes

	*/

	tokens := strings.Split(key, "/")
	tokens = append(tokens, "")
	copy(tokens[1+1:], tokens[1:])
	tokens[1] = "metadata"
	metadataPath := strings.Join(tokens, "/")
	tokens[1] = "data"
	dataPath := strings.Join(tokens, "/")

	logger.Info("listing from the vault")
	list, err := c.Logical().List(metadataPath)
	if err != nil {
		return err
	}
	if list == nil {
		// Check for leaf data
		logger.Warn("no list data found at this node")
		logger.Info("reading from the vault")
		read, err := c.Logical().Read(dataPath)
		if err != nil {
			return err
		}
		if read == nil {
			logger.Warn("no read data found at this node")
		} else {
			logger.Info("found read data at this node")
			err := node.setData(read.Data["data"])
			if err != nil {
				logger.Warn(err)
			}
		}
	} else {
		logger.Info("found list data at this node")
		spew.Dump(list.Data)

		for _, val := range list.Data {
			if slice, ok := val.([]interface{}); ok {
				for _, v := range slice {
					if name, ok := v.(string); ok {
						if strings.HasSuffix(name, "/") {
							var folder Folder
							folder.init()
							folder.setName(name)
							logger.Info("adding child folder")
							node.addChild(&folder)
							deepErr := Find(key+"/"+name, config, &folder, stack+1)
							if deepErr != nil {
								return deepErr
							}
						} else {
							var leaf Leaf
							leaf.setName(name)
							logger.Info("adding child leaf")
							node.addChild(&leaf)
							deepErr := Find(key+"/"+name, config, &leaf, stack+1)
							if deepErr != nil {
								return deepErr
							}
						}
					}
				}
			} else {
				logger.Errorf("not implemented for: %#v\n", val)
			}
		}
	}

	if stack == 0 {
		logger.Warn("data at root node")
		spew.Dump(node)
	}

	return
}
