package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"hypalynx.com/hypabase/internal/routes"
	_ "modernc.org/sqlite"
)

type config struct {
	Host string
	Port string
}

type server struct {
	config *config
	db     *sql.DB
	logger *log.Logger
	router string
}

func InitDB(url string) (*sql.DB, error) {
	return sql.Open("libsql", url)
}

func NewServer(db *sql.DB, logger *log.Logger, router string) http.Handler {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	routes.Setup(mux, logger)
	// handler = someMiddleware(handler)
	// handler = checkAuthHeaders(handler)
	return handler
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	db, err := InitDB("./local.sqlite3")
	config := &config{Host: "127.0.0.1", Port: "8080"}
	if err != nil {
		return err
	}
	srv := NewServer(db, log.New(os.Stdout, "hypabase", 0), "not-a-router")

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (credit: Alessandro Rosetti)
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server; %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}
