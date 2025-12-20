// Package docs provides embedded documentation access and rendering.
package docs

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

//go:embed content/*.md
var content embed.FS

// Doc represents a documentation file.
type Doc struct {
	// Name is the filename without extension (e.g., "README")
	Name string
	// Title is a human-readable title derived from the name
	Title string
	// Description is a brief description of the doc's purpose
	Description string
	// Path is the full path within the embed FS
	Path string
}

// docMeta contains metadata for known documentation files.
var docMeta = map[string]struct {
	title       string
	description string
}{
	"README":          {"Getting Started", "Project overview, features, and installation"},
	"TOOLS":           {"Tool Reference", "Complete reference for all 57 MCP tools"},
	"TROUBLESHOOTING": {"Troubleshooting", "Common issues and solutions"},
	"QUICKSTART":      {"Quick Start Guide", "Step-by-step setup instructions"},
	"BUILD":           {"Building from Source", "Build, test, and release instructions"},
}

// List returns all available documentation files.
func List() ([]Doc, error) {
	entries, err := fs.ReadDir(content, "content")
	if err != nil {
		return nil, fmt.Errorf("failed to read docs directory: %w", err)
	}

	docs := make([]Doc, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		doc := Doc{
			Name: name,
			Path: path.Join("content", entry.Name()),
		}

		// Apply metadata if known
		if meta, ok := docMeta[name]; ok {
			doc.Title = meta.title
			doc.Description = meta.description
		} else {
			// Generate title from name
			doc.Title = strings.ReplaceAll(name, "_", " ")
			doc.Description = ""
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// Get retrieves a documentation file by name.
// Name should be without the .md extension.
func Get(name string) (*Doc, error) {
	// Validate name to prevent path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return nil, fmt.Errorf("invalid doc name: %s", name)
	}

	docPath := path.Join("content", name+".md")
	_, err := fs.Stat(content, docPath)
	if err != nil {
		return nil, fmt.Errorf("doc not found: %s", name)
	}

	doc := &Doc{
		Name: name,
		Path: docPath,
	}

	if meta, ok := docMeta[name]; ok {
		doc.Title = meta.title
		doc.Description = meta.description
	} else {
		doc.Title = strings.ReplaceAll(name, "_", " ")
	}

	return doc, nil
}

// Content returns the raw markdown content of a document.
func Content(name string) (string, error) {
	doc, err := Get(name)
	if err != nil {
		return "", err
	}

	data, err := fs.ReadFile(content, doc.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read doc: %w", err)
	}

	return string(data), nil
}
