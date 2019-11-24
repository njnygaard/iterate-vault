module github.com/njnygaard/iterate-vault/cli

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/njnygaard/iterate-vault/vault-iterator v0.0.0
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/yaml.v2 v2.2.5
)

replace github.com/njnygaard/iterate-vault/vault-iterator => ../vault-iterator
