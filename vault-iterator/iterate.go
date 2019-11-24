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
	setData(data map[string]interface{}) (err error)
	getData()(data *map[string]interface{})
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

func (f *Folder) setData(data map[string]interface{}) (err error) {
	return errors.New("not implemented for folders")
}

func (f *Folder) getData()(data *map[string]interface{}) {
	return nil
}

func (f *Folder) addChild(node Node) {
	switch v := node.(type) {
	case *Folder:
		*f.childFolders = append(*f.childFolders, *v)
	case *Leaf:
		*f.childLeaves = append(*f.childLeaves, *v)
	default:
		return
	}
}

type Leaf struct {
	data *map[string]interface{}
	name string
}

func (l *Leaf) init() {
	var data = make(map[string]interface{})
	l.data = &data
}

func (l *Leaf) setName(name string) {
	l.name = name
}

func (l *Leaf) getChildren() (children *[]Node) {
	return
}

func (l *Leaf) setData(data map[string]interface{}) (err error) {
	// https://stackoverflow.com/a/38105687/1236359
	// You have to dereference the pointer to change 'what it points to'
	*l.data = data

	// In this example, I am passing the pointer to an address that is only scoped here.
	// This doesn't work.
	return nil
}

func (l *Leaf) getData()(data *map[string]interface{}) {
	return l.data
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

	tokens := strings.Split(key, "/")
	tokens = append(tokens, "")
	copy(tokens[1+1:], tokens[1:])
	tokens[1] = "metadata"
	metadataPath := strings.Join(tokens, "/")
	tokens[1] = "data"
	dataPath := strings.Join(tokens, "/")

	list, err := c.Logical().List(metadataPath)
	if err != nil {
		return err
	}
	if list == nil {
		read, err := c.Logical().Read(dataPath)
		if err != nil {
			return err
		}
		if read == nil {
			logger.Warn("no read data found at this node")
		} else {
			logger.Info("found leaf data")
			var data map[string]interface{}
			var ok bool
			if data, ok = read.Data["data"].(map[string]interface{}); ok {
				err := node.setData(data)
				if err != nil {
					logger.Warn(err)
				}
			}else{
				logger.Warn("couldn't pass type assertion")
			}
		}
	} else {
		logger.Info("found list data")
		for _, val := range list.Data {
			if slice, ok := val.([]interface{}); ok {
				for _, v := range slice {
					if name, ok := v.(string); ok {
						if strings.HasSuffix(name, "/") {
							var folder Folder
							folder.init()
							folder.setName(name)
							node.addChild(&folder)
							deepErr := Find(key+"/"+name, config, &folder, stack+1)
							if deepErr != nil {
								return deepErr
							}
						} else {
							var leaf Leaf
							leaf.init()
							leaf.setName(name)
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
