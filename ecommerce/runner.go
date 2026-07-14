//go:build ignore

package main

import (
	"bufio"
	"context"
	"flag" // To handle command-line flags
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings" // For suffix checking
	"sync"
	"syscall"
)

type Service struct {
	Name string
	Dir  string
}

func main() {
	// Command line flags definition
	protoFlag := flag.Bool("proto", false, "Generate all protobuf files")
	cleanProtoFlag := flag.Bool("clean-proto", false, "Remove all generated protobuf files")
	cleanDBFlag := flag.Bool("clean-db", false, "Remove all SQLite database files (.db)")
	cleanFlag := flag.Bool("clean", false, "Remove all compiled service binaries")
	flag.Parse()

	services := []Service{
		{Name: "AUTH", Dir: "services/auth-service"},
		{Name: "CART", Dir: "services/cart-service"},
		{Name: "CATALOG", Dir: "services/catalog-service"},
		{Name: "ORDER", Dir: "services/order-service"},
		{Name: "PAYMENT", Dir: "services/payment-service"},
		{Name: "WEB-SERVER", Dir: "web/server"},
	}

	// Trigger specific functions based on flags
	if *cleanProtoFlag {
		cleanProto()
		return
	}

	if *protoFlag {
		generateProto()
		return
	}

	if *cleanDBFlag {
		cleanDB()
		return
	}

	if *cleanFlag {
		cleanBinaries(services)
		return
	}

	// Default behavior: Start all services
	runAll(services)
}

// -----------------------------------------------------------------------------
// Core Actions
// -----------------------------------------------------------------------------

func runAll(services []Service) {
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
	binaryName := filepath.Base(svc.Dir)
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	absBinaryPath, err := filepath.Abs(filepath.Join(svc.Dir, binaryName))
	if err != nil {
		return fmt.Errorf("unable to determine the absolute path: %w", err)
	}

	fmt.Printf("[%s] Compiling...\n", svc.Name)
	buildCmd := exec.Command("go", "build", "-o", binaryName)
	buildCmd.Dir = svc.Dir
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("compilation error: %w", err)
	}

	cmd := exec.CommandContext(ctx, absBinaryPath)
	cmd.Dir = svc.Dir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to start binary: %w", err)
	}

	go streamLogs(svc.Name, stdout)
	go streamLogs(svc.Name, stderr)

	<-ctx.Done()

	defer os.Remove(absBinaryPath)

	return cmd.Wait()
}

// -----------------------------------------------------------------------------
// Utility Functions
// -----------------------------------------------------------------------------

// generateProto recompiles the .proto files (requires 'protoc' installed on the host)
func generateProto() {
	fmt.Println("[ORCHESTRATOR] Generating protobuf files...")
	protos := []string{
		"proto/auth/auth.proto",
		"proto/cart/cart.proto",
		"proto/catalog/catalog.proto",
		"proto/order/order.proto",
		"proto/payment/payment.proto",
	}

	for _, p := range protos {
		fmt.Printf("Compiling: %s\n", p)
		cmd := exec.Command("protoc", "--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error generating proto for %s: %v. Make sure 'protoc' is installed and added to your PATH.", p, err)
		}
	}
	fmt.Println("[ORCHESTRATOR] Protobuf generation completed successfully.")
}

// cleanProto recursively removes generated protobuf files (*.pb.go and *_grpc.pb.go)
func cleanProto() {
	fmt.Println("[ORCHESTRATOR] Removing generated protobuf files...")
	err := filepath.Walk("proto", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip folder if it doesn't exist yet
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".pb.go") || strings.HasSuffix(path, "_grpc.pb.go")) {
			fmt.Printf("Removing: %s\n", path)
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error cleaning protobuf files: %v", err)
	}
	fmt.Println("[ORCHESTRATOR] Protobuf cleanup completed.")
}

// cleanDB recursively scans the project directory and deletes all .db SQLite files
func cleanDB() {
	fmt.Println("[ORCHESTRATOR] Searching for SQLite database files...")
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".db" {
			fmt.Printf("Removing database: %s\n", path)
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error cleaning database files: %v", err)
	}
	fmt.Println("[ORCHESTRATOR] Database cleanup completed.")
}

// cleanBinaries deletes compiled binaries inside the service directories
func cleanBinaries(services []Service) {
	fmt.Println("[ORCHESTRATOR] Cleaning binaries...")
	for _, svc := range services {
		binaryName := filepath.Base(svc.Dir)
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(svc.Dir, binaryName)
		if _, err := os.Stat(binaryPath); err == nil {
			fmt.Printf("Removing binary: %s\n", binaryPath)
			os.Remove(binaryPath)
		}
	}
	fmt.Println("[ORCHESTRATOR] Binary cleanup completed.")
}

func streamLogs(prefix string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}
