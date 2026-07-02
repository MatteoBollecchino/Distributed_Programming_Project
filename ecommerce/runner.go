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
	"runtime" // <-- Cambiato: usiamo runtime per rilevare l'OS
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
		fmt.Println("\n[ORCHESTRATORE] Spegnimento di tutti i servizi in corso...")
		cancel()
	}()

	var wg sync.WaitGroup

	fmt.Println("[ORCHESTRATORE] Avvio della compilazione e dei servizi...")
	for _, svc := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := runService(ctx, s); err != nil {
				log.Printf("[%s] Errore: %v", s.Name, err)
			}
		}(svc)
	}

	wg.Wait()
	fmt.Println("[ORCHESTRATORE] Tutti i servizi sono stati arrestati.")
}

func runService(ctx context.Context, svc Service) error {
	// 1. Determina il nome del binario in base all'OS
	binaryName := filepath.Base(svc.Dir)
	if runtime.GOOS == "windows" { // <-- Corretto con runtime.GOOS
		binaryName += ".exe"
	}

	// 2. Ottieni il percorso ASSOLUTO del binario per evitare bug di pathing
	absBinaryPath, err := filepath.Abs(filepath.Join(svc.Dir, binaryName))
	if err != nil {
		return fmt.Errorf("impossibile determinare il percorso assoluto: %w", err)
	}

	// 3. Compilazione del servizio all'interno della sua cartella
	fmt.Printf("[%s] Compilazione in corso...\n", svc.Name)
	buildCmd := exec.Command("go", "build", "-o", binaryName)
	buildCmd.Dir = svc.Dir
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("errore di compilazione: %w", err)
	}

	// 4. Esecuzione del binario usando il percorso assoluto
	cmd := exec.CommandContext(ctx, absBinaryPath)
	cmd.Dir = svc.Dir // Mantiene la cartella di lavoro corretta per i database SQLite relativi

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("impossibile avviare il binario: %w", err)
	}

	go streamLogs(svc.Name, stdout)
	go streamLogs(svc.Name, stderr)

	<-ctx.Done()

	// Pulizia facoltativa dei binari generati alla chiusura
	defer os.Remove(absBinaryPath)

	return cmd.Wait()
}

func streamLogs(prefix string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}
