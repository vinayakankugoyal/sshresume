package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// TreeNode represents a file or directory in the tree structure.
type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*TreeNode
}

// Config represents the application configuration.
type Config struct {
	Folder string `yaml:"folder"`
	Tree   *TreeNode
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config.
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Build tree structure from folder.
	tree, err := buildTree(config.Folder)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}
	config.Tree = tree

	return &config, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Folder == "" {
		return fmt.Errorf("folder is required")
	}

	// Check if folder exists.
	info, err := os.Stat(c.Folder)
	if err != nil {
		return fmt.Errorf("folder does not exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("folder must be a directory")
	}

	return nil
}

// buildTree recursively builds a tree structure from a directory.
func buildTree(rootPath string) (*TreeNode, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Name:  filepath.Base(rootPath),
		Path:  rootPath,
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return nil, err
		}

		// Filter and sort entries.
		var children []*TreeNode
		for _, entry := range entries {
			// Skip hidden files.
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			childPath := filepath.Join(rootPath, entry.Name())
			childNode, err := buildTree(childPath)
			if err != nil {
				continue // Skip files that can't be read.
			}

			// Only include directories and markdown files.
			if childNode.IsDir || strings.HasSuffix(childNode.Name, ".md") {
				children = append(children, childNode)
			}
		}

		// Sort: directories first, then files, alphabetically.
		sort.Slice(children, func(i, j int) bool {
			if children[i].IsDir != children[j].IsDir {
				return children[i].IsDir
			}
			return children[i].Name < children[j].Name
		})

		node.Children = children
	}

	return node, nil
}
