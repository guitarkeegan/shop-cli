package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	"github.com/guitarkeegan/shop-cli/internal/examplesync"
)

var (
	ErrMissingEnvVars = errors.New("missing environment variables for examples sync")
	ErrTimeout        = errors.New("timed out syncing examples")
)

// Environment variables configuring "shop examples pull".
const (
	envGithubToken  = "SHOP_GITHUB_TOKEN"
	envGithubOwner  = "SHOP_GITHUB_OWNER"
	envGithubRepo   = "SHOP_GITHUB_REPO"
	envExamplesPath = "SHOP_EXAMPLES_PATH"
	envCacheDir     = "SHOP_CACHE_DIR"
)

type examplesPullOptions struct {
	Force bool
}

func newExamplesCmd() *clic.Clic {
	opts := &examplesPullOptions{}
	pull := clic.NewFromFunc(pullExamples(opts), "pull")
	pull.Flag(&opts.Force, "f|force", "Download all examples, even if unchanged.")

	return clic.NewFromFunc(requireSubcommand, "examples", pull)
}

// examplesConfigFromEnv reads the examples repo settings from the environment.
func examplesConfigFromEnv() (examplesync.Config, error) {
	cfg := examplesync.Config{
		Token:      os.Getenv(envGithubToken),
		Owner:      os.Getenv(envGithubOwner),
		Repo:       os.Getenv(envGithubRepo),
		FolderPath: os.Getenv(envExamplesPath),
	}
	if cfg.Token == "" || cfg.Owner == "" || cfg.Repo == "" || cfg.FolderPath == "" {
		return cfg, fmt.Errorf("%w: set %s, %s, %s and %s",
			ErrMissingEnvVars, envGithubToken, envGithubOwner, envGithubRepo, envExamplesPath)
	}
	return cfg, nil
}

func examplesDir() (string, error) {
	if dir := os.Getenv(envCacheDir); dir != "" {
		return filepath.Join(dir, "examples"), nil
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("find cache directory: %w", err)
	}
	return filepath.Join(cacheDir, "shop", "examples"), nil
}

func pullExamples(opts *examplesPullOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		cfg, err := examplesConfigFromEnv()
		if err != nil {
			return err
		}

		if cfg.DestDir, err = examplesDir(); err != nil {
			return err
		}
		if err := os.MkdirAll(cfg.DestDir, 0o755); err != nil {
			return fmt.Errorf("create cache directory: %w", err)
		}

		ctx, cancel := context.WithTimeoutCause(ctx, 30*time.Second, ErrTimeout)
		defer cancel()

		if err := examplesync.New(cfg).DownloadFolder(ctx, opts.Force); err != nil {
			return err
		}

		fmt.Println(colors.Success("Examples synced to " + cfg.DestDir))
		return nil
	}
}
