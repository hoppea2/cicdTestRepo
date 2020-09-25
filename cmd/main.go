package main

import (
	"cicdTestRepo"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8888", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our napodate service
	srv := cicdTestRepo.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := cicdTestRepo.Endpoints{
		GetEndpoint:      cicdTestRepo.MakeGetEndpoint(srv),
		StatusEndpoint:   cicdTestRepo.MakeStatusEndpoint(srv),
		ValidateEndpoint: cicdTestRepo.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("cicdTestRepo is listening on port:", *httpAddr)
		handler := cicdTestRepo.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
