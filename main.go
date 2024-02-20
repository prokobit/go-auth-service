package main

import (
	"flag"
	"fmt"

	"github.com/prokobit/auth-service/router"
)

type AppConfig struct {
	listenAddr string
}

func NewAppConfig(listenAddr string) *AppConfig {
	return &AppConfig{
		listenAddr: listenAddr,
	}
}

func main() {
	port := flag.Int("port", 8080, "Port")
	flag.Parse()
	listenAddr := fmt.Sprintf(":%d", *port)

	router.New(listenAddr).Start()
}
