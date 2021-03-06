package fetch

import (
	"context"
	"log"

	"github.com/hown3d/kevo/pkg/types"
)

type Fetcher interface {
	Fetch(context.Context, chan types.Image, chan error)
}

// Register can be called from init() on a plugin in this package
// It will automatically be added to the Fetchers map to be called externally
func Register(name string, f FetcherFactory) {
	fetcher, err := f()
	if err != nil {
		log.Printf("Error registering fetcher: %v", err)
		return
	}
	Fetchers[name] = fetcher
}

// InputFactory lets us use a closure to get intsances of the plugin struct
type FetcherFactory func() (Fetcher, error)

// Inputs registry
var Fetchers = map[string]Fetcher{}
