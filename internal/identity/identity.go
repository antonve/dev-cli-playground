package identity

import (
	"os"
	"time"
)

type Value struct {
	Service       string    `json:"service"`
	Branch        string    `json:"branch"`
	Revision      string    `json:"revision"`
	Pod           string    `json:"pod"`
	ObservedRoute string    `json:"observedRoute"`
	Time          time.Time `json:"time"`
}

func Current(service, route string) Value {
	branch := os.Getenv("DEV_BRANCH")
	if branch == "" {
		branch = "base"
	}
	revision := os.Getenv("DEV_REVISION")
	if revision == "" {
		revision = "base"
	}
	pod := os.Getenv("HOSTNAME")
	if pod == "" {
		pod = "local"
	}
	return Value{Service: service, Branch: branch, Revision: revision, Pod: pod, ObservedRoute: route, Time: time.Now().UTC()}
}
