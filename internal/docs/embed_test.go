package docs

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	docs, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(docs) == 0 {
		t.Error("List() returned empty list, expected at least some docs")
	}

	// Check expected docs are present
	expectedDocs := []string{"README", "TOOLS", "TROUBLESHOOTING", "QUICKSTART", "BUILD"}
	for _, expected := range expectedDocs {
		found := false
		for _, doc := range docs {
			if doc.Name == expected {
				found = true
				if doc.Title == "" {
					t.Errorf("Doc %s has empty title", expected)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected doc %s not found in list", expected)
		}
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		docName string
		wantErr bool
	}{
		{"valid doc", "README", false},
		{"valid doc with content", "TOOLS", false},
		{"invalid doc", "NONEXISTENT", true},
		{"path traversal attempt", "../../../etc/passwd", true},
		{"path with slash", "foo/bar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Get(tt.docName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get(%s) error = %v, wantErr %v", tt.docName, err, tt.wantErr)
			}
			if !tt.wantErr && doc == nil {
				t.Errorf("Get(%s) returned nil doc without error", tt.docName)
			}
		})
	}
}

func TestContent(t *testing.T) {
	content, err := Content("README")
	if err != nil {
		t.Fatalf("Content(README) error: %v", err)
	}

	if len(content) == 0 {
		t.Error("Content(README) returned empty content")
	}

	// README should contain project name
	if !strings.Contains(content, "agentic-obs") {
		t.Error("Content(README) should contain 'agentic-obs'")
	}
}

func TestContentInvalidDoc(t *testing.T) {
	_, err := Content("NONEXISTENT")
	if err == nil {
		t.Error("Content(NONEXISTENT) should return error")
	}
}
