// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// MIME type constants for UI resources.
const (
	MIMETypeHTML      = "text/html"
	MIMETypeURLList   = "text/uri-list"
	MIMETypeRemoteDOM = "application/vnd.mcp-ui.remote-dom"
)

// Framework constants for Remote DOM rendering.
type Framework string

const (
	// FrameworkReact specifies React as the rendering framework.
	FrameworkReact Framework = "react"
	// FrameworkWebComponents specifies Web Components as the rendering framework.
	FrameworkWebComponents Framework = "webcomponents"
)

// URIScheme is the URI scheme for UI resources.
const URIScheme = "ui://"

// Annotations contains metadata annotations for UI content.
// This mirrors the annotations concept from the MCP protocol.
type Annotations struct {
	// Audience specifies the intended audience for this content.
	Audience []string `json:"audience,omitempty"`
	// Priority indicates the relative importance of this content.
	Priority *float64 `json:"priority,omitempty"`
}

// UIContent is an [HTMLContent], [URLContent], or [RemoteDOMContent].
// This interface mirrors mcp.Content for UI resources.
type UIContent interface {
	// MarshalJSON serializes the content to JSON wire format.
	MarshalJSON() ([]byte, error)
	// mimeType returns the MIME type for this content.
	mimeType() string
	// fromWire populates the content from wire format.
	// Returns an error if the wire content cannot be parsed.
	fromWire(*wireUIContent) error
}

// HTMLContent contains inline HTML to render in a sandboxed iframe.
// The HTML is rendered using the iframe's srcdoc attribute.
//
// # Security
//
// This content is rendered in a sandboxed iframe with restricted permissions.
// However, the HTML is NOT sanitized by this SDK. Clients MUST ensure the
// iframe uses appropriate sandbox attributes (e.g., "allow-scripts" only when
// necessary) and implements Content Security Policy (CSP) headers. Server
// implementations should validate and sanitize HTML content before including
// it in responses.
type HTMLContent struct {
	// HTML is the inline HTML content to render.
	HTML string
	// Annotations contains optional metadata.
	Annotations *Annotations
}

// MarshalJSON serializes HTMLContent to the wire format.
func (c *HTMLContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireUIContent{
		MIMEType:    MIMETypeHTML,
		Text:        c.HTML,
		Annotations: c.Annotations,
	})
}

func (c *HTMLContent) mimeType() string { return MIMETypeHTML }

func (c *HTMLContent) fromWire(wire *wireUIContent) error {
	c.HTML = wire.Text
	c.Annotations = wire.Annotations
	return nil
}

// URLContent contains an external URL to render in an iframe.
// The URL is loaded using the iframe's src attribute.
type URLContent struct {
	// URL is the external URL to load.
	URL string
	// Annotations contains optional metadata.
	Annotations *Annotations
}

// Validate checks that the URLContent has a valid URL.
func (c *URLContent) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("URLContent URL is required")
	}
	parsed, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL must have http or https scheme, got: %s", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

// MarshalJSON serializes URLContent to the wire format.
func (c *URLContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireUIContent{
		MIMEType:    MIMETypeURLList,
		Text:        c.URL,
		Annotations: c.Annotations,
	})
}

func (c *URLContent) mimeType() string { return MIMETypeURLList }

func (c *URLContent) fromWire(wire *wireUIContent) error {
	c.URL = wire.Text
	c.Annotations = wire.Annotations
	return nil
}

// RemoteDOMContent contains a script for remote DOM rendering.
// The script is executed in a Web Worker inside a sandboxed iframe,
// and DOM changes are communicated to the host via JSON messages.
//
// # Security
//
// The Script field contains executable JavaScript code that runs in the client.
// While Remote DOM uses a Web Worker inside a sandboxed iframe for isolation,
// this SDK does NOT validate or sanitize the script content. Server implementations
// MUST ensure scripts are from trusted sources and clients SHOULD apply appropriate
// Content Security Policy (CSP) restrictions. Malicious scripts could attempt to
// exfiltrate data or abuse the action channel.
type RemoteDOMContent struct {
	// Script is the JavaScript code that constructs the remote DOM.
	Script string
	// Framework specifies the rendering framework (React or WebComponents).
	Framework Framework
	// Annotations contains optional metadata.
	Annotations *Annotations
}

// MarshalJSON serializes RemoteDOMContent to the wire format.
func (c *RemoteDOMContent) MarshalJSON() ([]byte, error) {
	mimeType := MIMETypeRemoteDOM + "+javascript"
	if c.Framework != "" {
		mimeType += "; framework=" + string(c.Framework)
	}
	return json.Marshal(&wireUIContent{
		MIMEType:    mimeType,
		Text:        c.Script,
		Annotations: c.Annotations,
	})
}

func (c *RemoteDOMContent) mimeType() string {
	mimeType := MIMETypeRemoteDOM + "+javascript"
	if c.Framework != "" {
		mimeType += "; framework=" + string(c.Framework)
	}
	return mimeType
}

func (c *RemoteDOMContent) fromWire(wire *wireUIContent) error {
	c.Script = wire.Text
	c.Annotations = wire.Annotations
	// Parse framework from MIME type (e.g., "application/vnd.mcp-ui.remote-dom+javascript; framework=react")
	if idx := strings.Index(wire.MIMEType, "framework="); idx != -1 {
		frameworkPart := wire.MIMEType[idx+len("framework="):]
		// Trim any trailing parameters (e.g., "; other=value")
		if endIdx := strings.IndexAny(frameworkPart, "; "); endIdx != -1 {
			frameworkPart = frameworkPart[:endIdx]
		}
		c.Framework = Framework(frameworkPart)
	}
	return nil
}

// BlobContent contains binary data (base64-encoded) for UI resources.
// This is used for images, fonts, or other binary assets.
type BlobContent struct {
	// Data is the binary content.
	Data []byte
	// MIMEType is the MIME type of the binary content.
	ContentMIMEType string
	// Annotations contains optional metadata.
	Annotations *Annotations
}

// MarshalJSON serializes BlobContent to the wire format.
func (c *BlobContent) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(c.Data)
	return json.Marshal(&wireUIContent{
		MIMEType:    c.ContentMIMEType,
		Blob:        encoded,
		Annotations: c.Annotations,
	})
}

func (c *BlobContent) mimeType() string { return c.ContentMIMEType }

func (c *BlobContent) fromWire(wire *wireUIContent) error {
	if wire.Blob != "" {
		data, err := base64.StdEncoding.DecodeString(wire.Blob)
		if err != nil {
			return fmt.Errorf("failed to decode base64 blob: %w", err)
		}
		c.Data = data
	}
	c.ContentMIMEType = wire.MIMEType
	c.Annotations = wire.Annotations
	return nil
}

// wireUIContent is the wire format for UI content.
// It represents all content types in a single structure for JSON marshaling.
type wireUIContent struct {
	MIMEType    string       `json:"mimeType"`
	Text        string       `json:"text,omitempty"`
	Blob        string       `json:"blob,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ContentFromWire converts wire format to the appropriate UIContent type.
func ContentFromWire(wire *wireUIContent) (UIContent, error) {
	if wire == nil {
		return nil, fmt.Errorf("nil wire content")
	}

	switch {
	case wire.MIMEType == MIMETypeHTML:
		c := &HTMLContent{}
		if err := c.fromWire(wire); err != nil {
			return nil, err
		}
		return c, nil
	case wire.MIMEType == MIMETypeURLList:
		c := &URLContent{}
		if err := c.fromWire(wire); err != nil {
			return nil, err
		}
		return c, nil
	case len(wire.MIMEType) >= len(MIMETypeRemoteDOM) && wire.MIMEType[:len(MIMETypeRemoteDOM)] == MIMETypeRemoteDOM:
		c := &RemoteDOMContent{}
		if err := c.fromWire(wire); err != nil {
			return nil, err
		}
		return c, nil
	case wire.Blob != "":
		c := &BlobContent{}
		if err := c.fromWire(wire); err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, fmt.Errorf("unknown content MIME type: %s", wire.MIMEType)
	}
}
