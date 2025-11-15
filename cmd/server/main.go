package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "be_test/internal/httpapi"
    "be_test/internal/store"
)

func main() {
    // Compose dependencies
    mem := store.NewMemoryItemStore()
    handler := httpapi.NewServer(mem)

    addr := ":8080"
    srv := &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Server listening on %s", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen error: %v", err)
        }
    }()

    // Graceful shutdown on SIGINT/SIGTERM
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop
    log.Println("Shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("server shutdown error: %v", err)
    }
    log.Println("Server stopped")
}