package main

import "os"

type Conf struct {
	Port   string `json:"port"`
	GWPort string `json:"gw_port"`
}

const (
	PortEnvKey   = "TYK_BROADCAST_PORT"
	GWPortEnvKey = "TYK_BROADCAST_GW_LISTENPORT"
)

func NewConfig() *Conf {
	p := os.Getenv(PortEnvKey)
	if p == "" {
		p = "8932"
	}

	gp := os.Getenv(GWPortEnvKey)
	if gp == "" {
		gp = "8080"
	}

	return &Conf{Port: p, GWPort: gp}
}
