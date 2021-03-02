package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/csothen/htmlparser/internal/controller"
	"github.com/gorilla/mux"
)

var (
	logger            *log.Logger                   = log.New(os.Stdout, "[ Parser ] - ", log.LstdFlags)
	parsingController *controller.ParsingController = controller.NewParsingController(logger)
)

func main() {
	fmt.Println("HTML Parser")

	router := configRouter()
	server := configServer(router)

	runServer(server)
}

func configRouter() *mux.Router {
	router := mux.NewRouter()

	postRouter := router.Methods(http.MethodPost).Subrouter()

	postRouter.HandleFunc("/api/parse", parsingController.ParseWebsite)

	return router
}

func configServer(router *mux.Router) *http.Server {
	server := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server
}

func runServer(server *http.Server) {
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	logger.Printf("Recieved terminate, graceful shutdown\nSignal recieved: %v", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
