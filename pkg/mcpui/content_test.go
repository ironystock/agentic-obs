// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLContent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		content  *HTMLContent
		wantJSON map[string]any
	}{
		{
			name: "simple HTML",
			content: &HTMLContent{
				HTML: "<div>Hello</div>",
			},
			wantJSON: map[string]any{
				"mimeType": "text/html",
				"text":     "<div>Hello</div>",
			},
		},
		{
			name: "HTML with annotations",
			content: &HTMLContent{
				HTML: "<p>Test</p>",
				Annotations: &Annotations{
					Audience: []string{"user"},
				},
			},
			wantJSON: map[string]any{
				"mimeType": "text/html",
				"text":     "<p>Test</p>",
				"annotations": map[string]any{
					"audience": []any{"user"},
				},
			},
		},
		{
			name: "empty HTML",
			content: &HTMLContent{
				HTML: "",
			},
			wantJSON: map[string]any{
				"mimeType": "text/html",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.content.MarshalJSON()
			require.NoError(t, err)

			var gotMap map[string]any
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)

			assert.Equal(t, tt.wantJSON["mimeType"], gotMap["mimeType"])
			if tt.wantJSON["text"] != nil {
				assert.Equal(t, tt.wantJSON["text"], gotMap["text"])
			}
		})
	}
}

func TestURLContent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		content  *URLContent
		wantJSON map[string]any
	}{
		{
			name: "simple URL",
			content: &URLContent{
				URL: "https://example.com",
			},
			wantJSON: map[string]any{
				"mimeType": "text/uri-list",
				"text":     "https://example.com",
			},
		},
		{
			name: "URL with path",
			content: &URLContent{
				URL: "https://example.com/dashboard?tab=settings",
			},
			wantJSON: map[string]any{
				"mimeType": "text/uri-list",
				"text":     "https://example.com/dashboard?tab=settings",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.content.MarshalJSON()
			require.NoError(t, err)

			var gotMap map[string]any
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)

			assert.Equal(t, tt.wantJSON["mimeType"], gotMap["mimeType"])
			assert.Equal(t, tt.wantJSON["text"], gotMap["text"])
		})
	}
}

func TestRemoteDOMContent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		content      *RemoteDOMContent
		wantMIMEType string
	}{
		{
			name: "without framework",
			content: &RemoteDOMContent{
				Script: "console.log('hello');",
			},
			wantMIMEType: "application/vnd.mcp-ui.remote-dom+javascript",
		},
		{
			name: "with React framework",
			content: &RemoteDOMContent{
				Script:    "React.createElement('div', null, 'Hello');",
				Framework: FrameworkReact,
			},
			wantMIMEType: "application/vnd.mcp-ui.remote-dom+javascript; framework=react",
		},
		{
			name: "with WebComponents framework",
			content: &RemoteDOMContent{
				Script:    "customElements.define('my-component', MyComponent);",
				Framework: FrameworkWebComponents,
			},
			wantMIMEType: "application/vnd.mcp-ui.remote-dom+javascript; framework=webcomponents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.content.MarshalJSON()
			require.NoError(t, err)

			var gotMap map[string]any
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)

			assert.Equal(t, tt.wantMIMEType, gotMap["mimeType"])
			assert.Equal(t, tt.content.Script, gotMap["text"])
		})
	}
}

func TestBlobContent_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		content    *BlobContent
		wantBlob   string
		blobOmited bool // true if blob should be omitted from JSON
	}{
		{
			name: "binary data",
			content: &BlobContent{
				Data:            []byte{0x89, 0x50, 0x4E, 0x47}, // PNG magic bytes
				ContentMIMEType: "image/png",
			},
			wantBlob: "iVBORw==", // base64 of the bytes
		},
		{
			name: "empty data",
			content: &BlobContent{
				Data:            []byte{},
				ContentMIMEType: "application/octet-stream",
			},
			blobOmited: true, // empty string is omitted by omitempty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.content.MarshalJSON()
			require.NoError(t, err)

			var gotMap map[string]any
			err = json.Unmarshal(got, &gotMap)
			require.NoError(t, err)

			assert.Equal(t, tt.content.ContentMIMEType, gotMap["mimeType"])
			if tt.blobOmited {
				_, exists := gotMap["blob"]
				assert.False(t, exists, "blob should be omitted for empty data")
			} else {
				assert.Equal(t, tt.wantBlob, gotMap["blob"])
			}
		})
	}
}

func TestContentFromWire(t *testing.T) {
	tests := []struct {
		name    string
		wire    *wireUIContent
		wantErr bool
		check   func(t *testing.T, c UIContent)
	}{
		{
			name: "HTML content",
			wire: &wireUIContent{
				MIMEType: MIMETypeHTML,
				Text:     "<div>Test</div>",
			},
			check: func(t *testing.T, c UIContent) {
				html, ok := c.(*HTMLContent)
				require.True(t, ok, "expected HTMLContent")
				assert.Equal(t, "<div>Test</div>", html.HTML)
			},
		},
		{
			name: "URL content",
			wire: &wireUIContent{
				MIMEType: MIMETypeURLList,
				Text:     "https://example.com",
			},
			check: func(t *testing.T, c UIContent) {
				url, ok := c.(*URLContent)
				require.True(t, ok, "expected URLContent")
				assert.Equal(t, "https://example.com", url.URL)
			},
		},
		{
			name: "RemoteDOM content",
			wire: &wireUIContent{
				MIMEType: MIMETypeRemoteDOM + "+javascript",
				Text:     "console.log('test');",
			},
			check: func(t *testing.T, c UIContent) {
				dom, ok := c.(*RemoteDOMContent)
				require.True(t, ok, "expected RemoteDOMContent")
				assert.Equal(t, "console.log('test');", dom.Script)
			},
		},
		{
			name: "RemoteDOM with framework",
			wire: &wireUIContent{
				MIMEType: MIMETypeRemoteDOM + "+javascript; framework=react",
				Text:     "React.render();",
			},
			check: func(t *testing.T, c UIContent) {
				dom, ok := c.(*RemoteDOMContent)
				require.True(t, ok, "expected RemoteDOMContent")
				assert.Equal(t, "React.render();", dom.Script)
			},
		},
		{
			name: "Blob content",
			wire: &wireUIContent{
				MIMEType: "image/png",
				Blob:     "iVBORw==",
			},
			check: func(t *testing.T, c UIContent) {
				blob, ok := c.(*BlobContent)
				require.True(t, ok, "expected BlobContent")
				assert.Equal(t, "image/png", blob.ContentMIMEType)
				assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47}, blob.Data)
			},
		},
		{
			name:    "nil wire",
			wire:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ContentFromWire(tt.wire)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestHTMLContent_MimeType(t *testing.T) {
	c := &HTMLContent{HTML: "<div>Test</div>"}
	assert.Equal(t, MIMETypeHTML, c.mimeType())
}

func TestURLContent_MimeType(t *testing.T) {
	c := &URLContent{URL: "https://example.com"}
	assert.Equal(t, MIMETypeURLList, c.mimeType())
}

func TestRemoteDOMContent_MimeType(t *testing.T) {
	tests := []struct {
		name      string
		framework Framework
		want      string
	}{
		{
			name:      "no framework",
			framework: "",
			want:      "application/vnd.mcp-ui.remote-dom+javascript",
		},
		{
			name:      "react",
			framework: FrameworkReact,
			want:      "application/vnd.mcp-ui.remote-dom+javascript; framework=react",
		},
		{
			name:      "webcomponents",
			framework: FrameworkWebComponents,
			want:      "application/vnd.mcp-ui.remote-dom+javascript; framework=webcomponents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RemoteDOMContent{Script: "test", Framework: tt.framework}
			assert.Equal(t, tt.want, c.mimeType())
		})
	}
}

func TestConstants(t *testing.T) {
	// Verify constants match MCP-UI specification
	assert.Equal(t, "text/html", MIMETypeHTML)
	assert.Equal(t, "text/uri-list", MIMETypeURLList)
	assert.Equal(t, "application/vnd.mcp-ui.remote-dom", MIMETypeRemoteDOM)
	assert.Equal(t, "ui://", URIScheme)
	assert.Equal(t, Framework("react"), FrameworkReact)
	assert.Equal(t, Framework("webcomponents"), FrameworkWebComponents)
}
