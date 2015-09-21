package adapter

import "net/http"

// Adapter (it gets its name from the adapter pattern — also known as the decorator pattern)
// is a function that both takes in and returns an http.Handler.
type Adapter func(http.Handler) http.Handler

// Adapt is a function that takes a handler you want to adapt, and a list of adapters. It
// iterates over all adapters, calling them one by one in a chained manner and returns a http.Handler again
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}
