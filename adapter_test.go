package adapter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func WithFirstHeader() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("first", "first")
			h.ServeHTTP(w, r)
		})
	}
}

func WithSecondHeader() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(w.Header().Get("first")) > 0 {
				w.Header().Add("second", "second") // only add second header when first is present (to verfiy order)
			}
			h.ServeHTTP(w, r)
		})
	}
}

func WithThirdHeader() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(w.Header().Get("second")) > 0 {
				w.Header().Add("third", "third") // only add third header when second is present (to verify order)
			}
			h.ServeHTTP(w, r)
		})
	}
}

func TestAdapt(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	helloHandler(w, req)
	a := Adapt(http.HandlerFunc(helloHandler), WithFirstHeader(), WithSecondHeader(), WithThirdHeader())
	a.ServeHTTP(w, req)
	// check if all headers are there, if not, order in Adapt function is wrong
	for _, header := range []string{"first", "second", "third"} {
		if len(w.Header().Get(header)) == 0 {
			t.Errorf("header %s not found", header)
		}
	}
	w = httptest.NewRecorder()
	a = Adapt(http.HandlerFunc(helloHandler), WithThirdHeader(), WithSecondHeader(), WithFirstHeader())
	a.ServeHTTP(w, req)
	// only first header can be present, otherwise order in Adapt function is wrong
	for _, header := range []string{"second", "third"} {
		if len(w.Header().Get(header)) != 0 {
			t.Errorf("header %s not found, but should not be there", header)
		}
	}
}
