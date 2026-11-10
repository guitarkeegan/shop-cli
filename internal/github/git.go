package git

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

type Config struct {
	token      string
	owner      string
	repo       string
	folderPath string
	destDir    string
}

type Client struct {
	config Config
	client *github.Client
}

type Manifest struct {
	LastSync time.Time         `json:"last_sync"`
	Files    map[string]string `json:"files"` // path -> SHA
}

func New(c Config) *Client {
	return &Client{
		config: c,
	}
}

func MakeConfig(ctx context.Context, token, owner, repo, folderPath, destDir string) Config {
	return Config{
		token:      token,
		owner:      owner,
		repo:       repo,
		folderPath: folderPath,
		destDir:    destDir,
	}
}

// will currently download all examples if any have changed
func (g *Client) DownloadFolder(ctx context.Context, force bool) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.config.token})
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
	if err := g.downloadFolderRecursive(ctx, g.config.folderPath, g.config.destDir, shas); err != nil {
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
	if err := g.collectRemoteSHAs(ctx, g.config.folderPath, remoteSHAs); err != nil {
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
	_, dirContents, _, err := g.client.Repositories.GetContents(ctx, g.config.owner, g.config.repo, path, nil)
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
	_, dirContents, _, err := g.client.Repositories.GetContents(ctx, g.config.owner, g.config.repo, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get contents: %w", err)
	}

	for _, content := range dirContents {
		destPath := filepath.Join(destDir, *content.Name)

		switch *content.Type {
		case "file":
			// Store SHA
			shas[*content.Path] = *content.SHA

			fileContent, _, _, err := g.client.Repositories.GetContents(ctx, g.config.owner, g.config.repo, *content.Path, nil)
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
	return filepath.Join(g.config.destDir, ".shop-manifest.json")
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
