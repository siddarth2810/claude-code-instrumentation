package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"claude-instrumentation/internal/hooks"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	//Create a context that automatically gets cancelled when Ctrl+C happens
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	httpServer := &http.Server{
		Addr:    ":10987",
		Handler: NewHandler(),
	}

	errs := make(chan error, 1)
	go func() {
		errs <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-errs:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	case <-ctx.Done(): //wait for Ctrl+c
	}

	//create shutdown ctx
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	return httpServer.Shutdown(shutdownCtx)
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", health)
	mux.HandleFunc("/hooks", handleHooks)
	return mux
}
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /health \n")
	io.WriteString(w, "Hello, from server!\n")
}

func handleHooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	input, err := hooks.DecodeHookInput(body)
	if err != nil {
		http.Error(w, "invalid hook input", http.StatusBadRequest)
		return
	}

	log.Printf("received hook event: %s", input.GetEventName())

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"recorded": true}); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
