//go:build ignore

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime" // we use runtime to detect the OS
	"sync"
	"syscall"
)

type Service struct {
	Name string
	Dir  string
}

func main() {
	services := []Service{
		{Name: "AUTH", Dir: "services/auth-service"},
		{Name: "CART", Dir: "services/cart-service"},
		{Name: "CATALOG", Dir: "services/catalog-service"},
		{Name: "ORDER", Dir: "services/order-service"},
		{Name: "PAYMENT", Dir: "services/payment-service"},
		{Name: "WEB-SERVER", Dir: "web/server"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n[ORCHESTRATOR] Shutting down all services...")
		cancel()
	}()

	var wg sync.WaitGroup

	fmt.Println("[ORCHESTRATOR] Starting compilation and services...")
	for _, svc := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := runService(ctx, s); err != nil {
				log.Printf("[%s] Error: %v", s.Name, err)
			}
		}(svc)
	}

	wg.Wait()
	fmt.Println("[ORCHESTRATOR] All services have been stopped.")
}

func runService(ctx context.Context, svc Service) error {
	// 1. Determine the binary name based on the OS
	binaryName := filepath.Base(svc.Dir)
	if runtime.GOOS == "windows" { // Corrected with runtime.GOOS
		binaryName += ".exe"
	}

	// 2. Get the ABSOLUTE path of the binary to avoid pathing bugs
	absBinaryPath, err := filepath.Abs(filepath.Join(svc.Dir, binaryName))
	if err != nil {
		return fmt.Errorf("unable to determine the absolute path: %w", err)
	}

	// 3. Service compilation inside its directory
	fmt.Printf("[%s] Compiling...\n", svc.Name)
	buildCmd := exec.Command("go", "build", "-o", binaryName)
	buildCmd.Dir = svc.Dir
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("compilation error: %w", err)
	}

	// 4. Execution of the binary using the absolute path
	cmd := exec.CommandContext(ctx, absBinaryPath)
	cmd.Dir = svc.Dir // Maintains the correct working directory for relative SQLite databases

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to start binary: %w", err)
	}

	go streamLogs(svc.Name, stdout)
	go streamLogs(svc.Name, stderr)

	<-ctx.Done()

	// Optional cleanup of generated binaries upon closure
	defer os.Remove(absBinaryPath)

	return cmd.Wait()
}

func streamLogs(prefix string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}
