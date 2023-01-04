package master

import "github.com/reindeer/talmud/pkg/master/provider"

func Get() string {
	return provider.Get()
}

func Save() {
	provider.Save()
}
