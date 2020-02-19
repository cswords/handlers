package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewProxyHandler creates a HandlerFunc func object which proxy an HTTP call by a key-value map as the configuration.
// The configuration should contain 2 fields: target, pathBase.
// The target field indicates the full URL as the proxy target.
// The request parameters after ? will be added to the parameters of the target URL.
// The pathBase field indicates the prefix of the full path of the URL of the original request received, which includes the prefix by router configuration and path by handler configuration.
func NewProxyHandler(config map[string]string) http.HandlerFunc {
	target := config["target"]
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	pathBase := config["pathBase"]

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, pathBase) {
			err := fmt.Errorf("Request URL %q does not match path base %q", r.URL.String(), pathBase)
			panic(err)
		}

		if r.Method != http.MethodOptions {

			log.Println("Request url is ", r.URL.String())
			log.Println("Target url is ", targetURL.String())

			r.URL.Scheme = targetURL.Scheme
			r.URL.Host = targetURL.Host
			r.URL.Path = targetURL.Path + r.URL.Path[len(pathBase):len(r.URL.Path)]
			r.URL.RawQuery = strings.Join([]string{targetURL.RawQuery, r.URL.RawQuery}, "&")

			r.Host = r.URL.Host
			r.RequestURI = r.URL.RequestURI()

			log.Println("Request url rewriten to ", r.URL.String())

			// request will be copied
			reverseProxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
