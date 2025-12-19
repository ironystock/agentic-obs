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

func TestUIResource_Validate(t *testing.T) {
	tests := []struct {
		name     string
		resource *UIResource
		wantErr  string
	}{
		{
			name: "valid resource",
			resource: &UIResource{
				URI:  "ui://dashboard/main",
				Name: "Dashboard",
			},
			wantErr: "",
		},
		{
			name: "missing URI",
			resource: &UIResource{
				Name: "Dashboard",
			},
			wantErr: "missing URI",
		},
		{
			name: "invalid URI scheme",
			resource: &UIResource{
				URI:  "http://example.com",
				Name: "Dashboard",
			},
			wantErr: "must start with ui://",
		},
		{
			name: "missing Name",
			resource: &UIResource{
				URI: "ui://dashboard/main",
			},
			wantErr: "missing Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resource.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestUIResourceContents_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		contents *UIResourceContents
		wantErr  bool
		check    func(t *testing.T, data []byte)
	}{
		{
			name: "text content",
			contents: &UIResourceContents{
				URI:      "ui://greeting/hello",
				MIMEType: MIMETypeHTML,
				Text:     "<div>Hello</div>",
			},
			check: func(t *testing.T, data []byte) {
				var m map[string]any
				require.NoError(t, json.Unmarshal(data, &m))
				assert.Equal(t, "ui://greeting/hello", m["uri"])
				assert.Equal(t, MIMETypeHTML, m["mimeType"])
				assert.Equal(t, "<div>Hello</div>", m["text"])
				_, hasBlob := m["blob"]
				assert.False(t, hasBlob)
			},
		},
		{
			name: "blob content",
			contents: &UIResourceContents{
				URI:      "ui://image/logo",
				MIMEType: "image/png",
				Blob:     []byte{0x89, 0x50, 0x4E, 0x47},
			},
			check: func(t *testing.T, data []byte) {
				var m map[string]any
				require.NoError(t, json.Unmarshal(data, &m))
				assert.Equal(t, "ui://image/logo", m["uri"])
				assert.Equal(t, "image/png", m["mimeType"])
				// blob should be base64 encoded
				assert.NotNil(t, m["blob"])
			},
		},
		{
			name: "empty blob",
			contents: &UIResourceContents{
				URI:      "ui://empty/data",
				MIMEType: "application/octet-stream",
				Blob:     []byte{},
			},
			check: func(t *testing.T, data []byte) {
				var m map[string]any
				require.NoError(t, json.Unmarshal(data, &m))
				// empty blob should still include the blob field
				_, hasBlob := m["blob"]
				assert.True(t, hasBlob, "blob field should be present for empty slice")
			},
		},
		{
			name: "missing URI",
			contents: &UIResourceContents{
				MIMEType: MIMETypeHTML,
				Text:     "<div>No URI</div>",
			},
			wantErr: true,
		},
		{
			name: "both text and blob",
			contents: &UIResourceContents{
				URI:      "ui://invalid/both",
				MIMEType: MIMETypeHTML,
				Text:     "text",
				Blob:     []byte{1, 2, 3},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.contents.MarshalJSON()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, data)
		})
	}
}

func TestNewUIResourceContents(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		content UIContent
		wantErr bool
		check   func(t *testing.T, rc *UIResourceContents)
	}{
		{
			name: "from HTMLContent",
			uri:  "ui://greeting/hello",
			content: &HTMLContent{
				HTML: "<div>Hello World</div>",
			},
			check: func(t *testing.T, rc *UIResourceContents) {
				assert.Equal(t, "ui://greeting/hello", rc.URI)
				assert.Equal(t, MIMETypeHTML, rc.MIMEType)
				assert.Equal(t, "<div>Hello World</div>", rc.Text)
			},
		},
		{
			name: "from URLContent",
			uri:  "ui://external/dashboard",
			content: &URLContent{
				URL: "https://example.com/dashboard",
			},
			check: func(t *testing.T, rc *UIResourceContents) {
				assert.Equal(t, "ui://external/dashboard", rc.URI)
				assert.Equal(t, MIMETypeURLList, rc.MIMEType)
				assert.Equal(t, "https://example.com/dashboard", rc.Text)
			},
		},
		{
			name: "from RemoteDOMContent",
			uri:  "ui://component/widget",
			content: &RemoteDOMContent{
				Script:    "console.log('hello');",
				Framework: FrameworkReact,
			},
			check: func(t *testing.T, rc *UIResourceContents) {
				assert.Equal(t, "ui://component/widget", rc.URI)
				assert.Contains(t, rc.MIMEType, MIMETypeRemoteDOM)
				assert.Contains(t, rc.MIMEType, "react")
				assert.Equal(t, "console.log('hello');", rc.Text)
			},
		},
		{
			name:    "empty URI",
			uri:     "",
			content: &HTMLContent{HTML: "test"},
			wantErr: true,
		},
		{
			name:    "nil content",
			uri:     "ui://test/nil",
			content: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := NewUIResourceContents(tt.uri, tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, rc)
		})
	}
}

func TestUIResourceContents_ToUIContent(t *testing.T) {
	tests := []struct {
		name     string
		contents *UIResourceContents
		check    func(t *testing.T, c UIContent)
	}{
		{
			name: "HTML content",
			contents: &UIResourceContents{
				URI:      "ui://test/html",
				MIMEType: MIMETypeHTML,
				Text:     "<p>Test</p>",
			},
			check: func(t *testing.T, c UIContent) {
				html, ok := c.(*HTMLContent)
				require.True(t, ok, "expected HTMLContent")
				assert.Equal(t, "<p>Test</p>", html.HTML)
			},
		},
		{
			name: "URL content",
			contents: &UIResourceContents{
				URI:      "ui://test/url",
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
			contents: &UIResourceContents{
				URI:      "ui://test/dom",
				MIMEType: MIMETypeRemoteDOM + "+javascript",
				Text:     "document.write('hello');",
			},
			check: func(t *testing.T, c UIContent) {
				dom, ok := c.(*RemoteDOMContent)
				require.True(t, ok, "expected RemoteDOMContent")
				assert.Equal(t, "document.write('hello');", dom.Script)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := tt.contents.ToUIContent()
			require.NoError(t, err)
			tt.check(t, content)
		})
	}
}

func TestUIResourceTemplate_Validate(t *testing.T) {
	tests := []struct {
		name     string
		template *UIResourceTemplate
		wantErr  string
	}{
		{
			name: "valid template",
			template: &UIResourceTemplate{
				URITemplate: "ui://dashboard/{id}",
				Name:        "Dashboard",
			},
			wantErr: "",
		},
		{
			name: "missing URITemplate",
			template: &UIResourceTemplate{
				Name: "Dashboard",
			},
			wantErr: "missing URITemplate",
		},
		{
			name: "invalid URI scheme",
			template: &UIResourceTemplate{
				URITemplate: "http://example.com/{id}",
				Name:        "Dashboard",
			},
			wantErr: "must start with ui://",
		},
		{
			name: "missing Name",
			template: &UIResourceTemplate{
				URITemplate: "ui://dashboard/{id}",
			},
			wantErr: "missing Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestUIResource_JSONSerialization(t *testing.T) {
	resource := &UIResource{
		URI:         "ui://dashboard/main",
		Name:        "main-dashboard",
		Title:       "Main Dashboard",
		Description: "The primary dashboard view",
		MIMEType:    MIMETypeHTML,
		Annotations: &Annotations{
			Audience: []string{"user"},
		},
	}

	data, err := json.Marshal(resource)
	require.NoError(t, err)

	var decoded UIResource
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resource.URI, decoded.URI)
	assert.Equal(t, resource.Name, decoded.Name)
	assert.Equal(t, resource.Title, decoded.Title)
	assert.Equal(t, resource.Description, decoded.Description)
	assert.Equal(t, resource.MIMEType, decoded.MIMEType)
	assert.Equal(t, resource.Annotations.Audience, decoded.Annotations.Audience)
}

func TestListUIResourcesResult_JSONSerialization(t *testing.T) {
	result := &ListUIResourcesResult{
		Resources: []*UIResource{
			{URI: "ui://test/one", Name: "one"},
			{URI: "ui://test/two", Name: "two"},
		},
		NextCursor: "cursor123",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded ListUIResourcesResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded.Resources, 2)
	assert.Equal(t, "cursor123", decoded.NextCursor)
}

func TestReadUIResourceResult_JSONSerialization(t *testing.T) {
	result := &ReadUIResourceResult{
		Contents: []*UIResourceContents{
			{
				URI:      "ui://test/content",
				MIMEType: MIMETypeHTML,
				Text:     "<div>Content</div>",
			},
		},
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded ReadUIResourceResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded.Contents, 1)
	assert.Equal(t, "ui://test/content", decoded.Contents[0].URI)
}
