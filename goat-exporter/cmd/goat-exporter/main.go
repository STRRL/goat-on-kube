package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"port"`
}

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

			slog.Info("starting goat-exporter", slog.String("port", cfg.Port))
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			<-ctx.Done()
			slog.Info("shutting down goat-exporter")
			return nil
		},
	}

	cmd.PersistentFlags().String("port", "8080", "port for the exporter to listen on")

	return cmd
}

func initConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("port", "8080")
}

func loadConfig() (Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}
	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}
