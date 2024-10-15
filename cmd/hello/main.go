package main

import (
	"context"
	"flag"

	"github.com/go-obvious/server"

	"github.com/go-obvious/server-example/internal/build"
	"github.com/go-obvious/server-example/internal/service/hello"
)

func parseFlags() (mode *string, port *uint, domain *string) {
	mode = flag.String("mode", "http", "Mode to run the server in (http or lambda)")
	port = flag.Uint("port", 8080, "Port to run the server on")
	domain = flag.String("domain", "metrics.cloudzero.com", "Domain for the server")
	flag.Parse()
	return
}

func main() {
	mode, port, domain := parseFlags()

	server.New(
		&server.Config{
			Domain: *domain,
			Port:   *port,
			Mode:   *mode,
		},
		&server.ServerVersion{
			Revision: build.Rev,
			Tag:      build.Tag,
			Time:     build.Time,
		},
		hello.NewService("/"),
	).Run(context.Background())
}
