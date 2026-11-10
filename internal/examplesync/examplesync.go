// Package examplesync downloads a folder of example files from a GitHub
// repo, skipping the download when nothing has changed.
package examplesync

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v82/github"
	"golang.org/x/oauth2"
)

// Config identifies the GitHub folder to sync and where to save it.
type Config struct {
	Token      string // GitHub token with read access to the repo
	Owner      string
	Repo       string
	FolderPath string // folder within the repo, e.g. "examples"
	DestDir    string // local directory to download into
}

// Client syncs a folder of examples from a GitHub repo to a local directory.
type Client struct {
	config Config
	client *github.Client
}

// Manifest records the SHA of each synced file, so unchanged folders can
// skip downloading.
type Manifest struct {
	LastSync time.Time         `json:"last_sync"`
	Files    map[string]string `json:"files"` // path -> SHA
}

func New(c Config) *Client {
	return &Client{config: c}
}

// DownloadFolder downloads every file in the folder if any file changed
// since the last sync, or always when force is true.
func (g *Client) DownloadFolder(ctx context.Context, force bool) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.config.Token})
	tc := oauth2.NewClient(ctx, ts)
	g.client = github.NewClient(tc)

	if !force {
		needsUpdate, err := g.needsUpdate(ctx)
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		if !needsUpdate {
			fmt.Println("✓ Examples are already up to date")
			return nil
		}
	}

	fmt.Println("Downloading examples...")
	shas := make(map[string]string)
	if err := g.downloadFolderRecursive(ctx, g.config.FolderPath, g.config.DestDir, shas); err != nil {
		return err
	}

	return g.saveManifest(shas)
}

func (g *Client) needsUpdate(ctx context.Context) (bool, error) {
	manifest, err := g.loadManifest()
	if err != nil {
		return true, nil
	}

	remoteSHAs := make(map[string]string)
	if err := g.collectRemoteSHAs(ctx, g.config.FolderPath, remoteSHAs); err != nil {
		return false, err
	}

	if len(remoteSHAs) != len(manifest.Files) {
		return true, nil
	}

	for path, remoteSHA := range remoteSHAs {
		if localSHA, exists := manifest.Files[path]; !exists || localSHA != remoteSHA {
			return true, nil
		}
	}

	return false, nil
}

// this is a recursive call, switch to git trees api if we end up with a lot of folders
func (g *Client) collectRemoteSHAs(ctx context.Context, path string, shas map[string]string) error {
	_, dirContents, _, err := g.client.Repositories.GetContents(ctx, g.config.Owner, g.config.Repo, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get contents of %s: %w", path, err)
	}

	for _, content := range dirContents {
		switch *content.Type {
		case "file":
			shas[*content.Path] = *content.SHA
		case "dir":
			if err := g.collectRemoteSHAs(ctx, *content.Path, shas); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Client) downloadFolderRecursive(ctx context.Context, path, destDir string, shas map[string]string) error {
	_, dirContents, _, err := g.client.Repositories.GetContents(ctx, g.config.Owner, g.config.Repo, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get contents: %w", err)
	}

	for _, content := range dirContents {
		destPath := filepath.Join(destDir, *content.Name)

		switch *content.Type {
		case "file":
			// Store SHA
			shas[*content.Path] = *content.SHA

			fileContent, _, _, err := g.client.Repositories.GetContents(ctx, g.config.Owner, g.config.Repo, *content.Path, nil)
			if err != nil {
				return fmt.Errorf("failed to get file %s: %w", *content.Path, err)
			}

			decoded, err := base64.StdEncoding.DecodeString(*fileContent.Content)
			if err != nil {
				return fmt.Errorf("failed to decode %s: %w", *content.Path, err)
			}

			if err := os.WriteFile(destPath, decoded, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", destPath, err)
			}

			fmt.Printf("  ✓ %s\n", *content.Path)

		case "dir":
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create dir %s: %w", destPath, err)
			}
			if err := g.downloadFolderRecursive(ctx, *content.Path, destPath, shas); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Client) manifestPath() string {
	return filepath.Join(g.config.DestDir, ".shop-manifest.json")
}

func (g *Client) loadManifest() (*Manifest, error) {
	data, err := os.ReadFile(g.manifestPath())
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (g *Client) saveManifest(shas map[string]string) error {
	manifest := Manifest{
		LastSync: time.Now(),
		Files:    shas,
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(g.manifestPath(), data, 0644)
}
