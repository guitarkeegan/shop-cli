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
	git "github.com/guitarkeegan/shop-cli/internal/github"
)

type ExampleOptions struct {
	Force bool
}

func NewExampleOptions() *ExampleOptions {
	return &ExampleOptions{}
}

var (
	ErrMissingEnvVars = errors.New("missing env vars for git usage")
	ErrTimeout        = errors.New("took too long to complete")
)

const (
	githubToken      = "SHOP_GITHUB_TOKEN"
	shopGithubOwner  = "SHOP_GITHUB_OWNER"
	shopGithubRepo   = "SHOP_GITHUB_REPO"
	shopExamplesPath = "SHOP_EXAMPLES_PATH"
	shopCacheDir     = "SHOP_CACHE_DIR"
)

func examplesHandler() clic.HandlerFunc {
	return func(ctx context.Context) error {
		return ErrMissingFlagOrCommand
	}
}

func hasRequiredEnvs() (string, string, string, string, error) {
	gitToken := os.Getenv(githubToken)
	repo := os.Getenv(shopGithubRepo)
	owner := os.Getenv(shopGithubOwner)
	githubPath := os.Getenv(shopExamplesPath) // GitHub folder path, e.g., "examples"

	if gitToken == "" || repo == "" || owner == "" || githubPath == "" {
		return "", "", "", "", fmt.Errorf("%w: need: %s, %s, %s, %s\n",
			ErrMissingEnvVars, githubToken, shopGithubOwner, shopGithubRepo, shopExamplesPath)
	}
	return gitToken, repo, owner, githubPath, nil
}

func getExamplesDir() (string, error) {
	if dir := os.Getenv(shopCacheDir); dir != "" {
		return filepath.Join(dir, "examples"), nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "shop", "examples"), nil
}

func handleExamples(opt *ExampleOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		gitToken, repo, owner, githubPath, err := hasRequiredEnvs()
		if err != nil {
			return err
		}

		localDir, err := getExamplesDir()
		if err != nil {
			return fmt.Errorf("getExamplesDir: %s", err)
		}

		if err := os.MkdirAll(localDir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory: %w", err)
		}

		timedCtx, cancel := context.WithTimeoutCause(ctx, time.Second*30, ErrTimeout)
		defer cancel()

		cfg := git.MakeConfig(ctx, gitToken, owner, repo, githubPath, localDir)
		client := git.New(cfg)

		if err := client.DownloadFolder(timedCtx, opt.Force); err != nil {
			return err
		}

		colors.Success(fmt.Sprintf("Examples synced to: %s", localDir))
		return nil
	}
}
