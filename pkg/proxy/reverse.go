package proxy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func createLocalReverseDirector(portMappings map[string]int) func(*http.Request) {
	return func(r *http.Request) {
		domain := getDomainFromHost(r.Host)
		log.Printf("director sees %s\n", domain)
		targetPort, found := portMappings[domain]
		if !found {
			log.Printf("unknown proxy domain %s\n", domain)
			return
		}

		r.URL.Host = fmt.Sprintf("localhost:%d", targetPort)
		r.URL.Scheme = "http"
	}
}

func getDomainFromHost(host string) string {
	colonPosition := strings.Index(host, ":")
	if colonPosition == -1 {
		return host
	}
	return host[0:colonPosition]
}
