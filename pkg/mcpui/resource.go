// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// UIResource represents an interactive UI resource definition.
// This mirrors mcp.Resource for UI-specific resources.
// See https://mcpui.dev/guide/protocol-details
type UIResource struct {
	// URI is the unique identifier for this resource (e.g., "ui://dashboard/main").
	URI string `json:"uri"`
	// Name is intended for programmatic or logical use.
	Name string `json:"name"`
	// Title is the human-readable display name.
	Title string `json:"title,omitempty"`
	// Description explains what this resource represents.
	Description string `json:"description,omitempty"`
	// MIMEType is the MIME type of this resource, if known.
	MIMEType string `json:"mimeType,omitempty"`
	// Annotations contains optional metadata.
	Annotations *Annotations `json:"annotations,omitempty"`
}

// Validate checks that the UIResource has required fields and valid URI scheme.
func (r *UIResource) Validate() error {
	if r.URI == "" {
		return errors.New("UIResource missing URI")
	}
	if !strings.HasPrefix(r.URI, URIScheme) {
		return errors.New("UIResource URI must start with " + URIScheme)
	}
	if r.Name == "" {
		return errors.New("UIResource missing Name")
	}
	return nil
}

// UIResourceContents contains the contents of a specific UI resource.
// This mirrors mcp.ResourceContents for UI-specific resources.
type UIResourceContents struct {
	// URI is the resource identifier.
	URI string `json:"uri"`
	// MIMEType is the content MIME type.
	MIMEType string `json:"mimeType,omitempty"`
	// Text is the textual content (HTML, URL, or script).
	Text string `json:"text,omitempty"`
	// Blob is the binary content (base64-encoded in JSON).
	Blob []byte `json:"blob,omitempty"`
	// Annotations contains optional metadata.
	Annotations *Annotations `json:"annotations,omitempty"`
}

// MarshalJSON serializes UIResourceContents to JSON.
// It follows the mcp.ResourceContents pattern.
func (r *UIResourceContents) MarshalJSON() ([]byte, error) {
	if r.URI == "" {
		return nil, errors.New("UIResourceContents missing URI")
	}
	if r.Blob == nil {
		// Text content. Marshal normally.
		type wireResourceContents UIResourceContents // lacks MarshalJSON method
		return json.Marshal((wireResourceContents)(*r))
	}
	// Blob content.
	if r.Text != "" {
		return nil, errors.New("UIResourceContents has non-zero Text and Blob fields")
	}
	// r.Blob may be the empty slice, so marshal with an alternative definition
	// to ensure "blob" is always included.
	br := struct {
		URI         string       `json:"uri,omitempty"`
		MIMEType    string       `json:"mimeType,omitempty"`
		Blob        []byte       `json:"blob"`
		Annotations *Annotations `json:"annotations,omitempty"`
	}{
		URI:         r.URI,
		MIMEType:    r.MIMEType,
		Blob:        r.Blob,
		Annotations: r.Annotations,
	}
	return json.Marshal(br)
}

// NewUIResourceContents creates UIResourceContents from a UIContent.
func NewUIResourceContents(uri string, content UIContent) (*UIResourceContents, error) {
	if uri == "" {
		return nil, errors.New("URI is required")
	}
	if content == nil {
		return nil, errors.New("content is required")
	}

	// Marshal content to wire format to extract fields
	data, err := content.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var wire wireUIContent
	if err := json.Unmarshal(data, &wire); err != nil {
		return nil, err
	}

	rc := &UIResourceContents{
		URI:         uri,
		MIMEType:    wire.MIMEType,
		Annotations: wire.Annotations,
	}

	if wire.Blob != "" {
		// Decode base64 blob back to bytes
		decoded, err := base64.StdEncoding.DecodeString(wire.Blob)
		if err != nil {
			return nil, errors.New("failed to decode base64 blob: " + err.Error())
		}
		rc.Blob = decoded
	} else {
		rc.Text = wire.Text
	}

	return rc, nil
}

// ToUIContent converts UIResourceContents back to a UIContent.
func (r *UIResourceContents) ToUIContent() (UIContent, error) {
	wire := &wireUIContent{
		MIMEType:    r.MIMEType,
		Text:        r.Text,
		Annotations: r.Annotations,
	}
	if r.Blob != nil {
		// Encode blob to base64 string for wire format
		wire.Blob = base64.StdEncoding.EncodeToString(r.Blob)
	}
	return ContentFromWire(wire)
}

// UIResourceTemplate describes a template for UI resources.
// This mirrors mcp.ResourceTemplate for UI-specific resources.
type UIResourceTemplate struct {
	// URITemplate is a URI template following RFC 6570.
	URITemplate string `json:"uriTemplate"`
	// Name is intended for programmatic or logical use.
	Name string `json:"name"`
	// Title is the human-readable display name.
	Title string `json:"title,omitempty"`
	// Description explains what resources this template represents.
	Description string `json:"description,omitempty"`
	// MIMEType is the MIME type of resources matching this template.
	MIMEType string `json:"mimeType,omitempty"`
	// Annotations contains optional metadata.
	Annotations *Annotations `json:"annotations,omitempty"`
}

// Validate checks that the UIResourceTemplate has required fields.
func (t *UIResourceTemplate) Validate() error {
	if t.URITemplate == "" {
		return errors.New("UIResourceTemplate missing URITemplate")
	}
	if !strings.HasPrefix(t.URITemplate, URIScheme) {
		return errors.New("UIResourceTemplate URITemplate must start with " + URIScheme)
	}
	if t.Name == "" {
		return errors.New("UIResourceTemplate missing Name")
	}
	return nil
}

// ReadUIResourceResult is the result of reading a UI resource.
// This mirrors mcp.ReadResourceResult.
type ReadUIResourceResult struct {
	// Contents contains the resource contents.
	Contents []*UIResourceContents `json:"contents"`
}

// ListUIResourcesResult is the result of listing UI resources.
// This mirrors mcp.ListResourcesResult.
type ListUIResourcesResult struct {
	// Resources is the list of available UI resources.
	Resources []*UIResource `json:"resources"`
	// NextCursor is an opaque token for pagination.
	NextCursor string `json:"nextCursor,omitempty"`
}
