package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Service struct {
	Name string
	Dir  string
}

func main() {
	services := []Service{
		{"auth", "services/auth-service"},
		{"cart", "services/cart-service"},
		{"catalog", "services/catalog-service"},
		{"order", "services/order-service"},
		{"payment", "services/payment-service"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Intercetta Ctrl+C per fermare tutti
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Avvia tutti i servizi
	for _, s := range services {
		cmd := exec.CommandContext(ctx, "go", "run", "main.go")
		cmd.Dir = s.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		go func(name string, c *exec.Cmd) {
			if err := c.Run(); err != nil {
				log.Printf("[%s] stopped: %v", name, err)
			}
		}(s.Name, cmd)
	}

	<-sig
	log.Println("Stopping all services...")
	cancel()
}
