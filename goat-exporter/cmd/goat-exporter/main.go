package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/strrl/goat-on-kube/goat-exporter/internal/collector"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rootCmd := newRootCommand()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cobra.OnInitialize(initConfig)

	cmd := &cobra.Command{
		Use:   "goat-exporter",
		Short: "Run the Goat exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()

			if err != nil {
				return err
			}
			return run(cmd.Context(), cfg)
		},
	}

	cmd.PersistentFlags().String("port", "8080", "port for the exporter to listen on")

	return cmd
}

func run(ctx context.Context, cfg Config) error {
	slog.Info(
		"starting goat-exporter",
		slog.String("port", cfg.Port),
		slog.String("rpc_node", cfg.RPCNode),
	)

	reg := prometheus.NewRegistry()

	rpcCollector, err := collector.NewGethRPCCollector(cfg.RPCNode)
	if err != nil {
		return fmt.Errorf("init collector: %w", err)
	}
	defer rpcCollector.Close()

	if err := reg.Register(rpcCollector); err != nil {
		return fmt.Errorf("register collector: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	errCh := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	slog.Info("shutting down goat-exporter")
	return nil
}
