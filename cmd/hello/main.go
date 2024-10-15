package main

import (
	"context"

	"github.com/go-obvious/server"

	"github.com/go-obvious/server-example/internal/build"
	"github.com/go-obvious/server-example/internal/service/hello"
)

func main() {
	server.New(
		&server.ServerVersion{
			Revision: build.Rev,
			Tag:      build.Tag,
			Time:     build.Time,
		},
		hello.NewService("/"),
	).Run(context.Background())
}
