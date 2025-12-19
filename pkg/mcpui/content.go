// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	fromWire(*wireUIContent)
}

// HTMLContent contains inline HTML to render in a sandboxed iframe.
// The HTML is rendered using the iframe's srcdoc attribute.
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

func (c *HTMLContent) fromWire(wire *wireUIContent) {
	c.HTML = wire.Text
	c.Annotations = wire.Annotations
}

// URLContent contains an external URL to render in an iframe.
// The URL is loaded using the iframe's src attribute.
type URLContent struct {
	// URL is the external URL to load.
	URL string
	// Annotations contains optional metadata.
	Annotations *Annotations
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

func (c *URLContent) fromWire(wire *wireUIContent) {
	c.URL = wire.Text
	c.Annotations = wire.Annotations
}

// RemoteDOMContent contains a script for remote DOM rendering.
// The script is executed in a Web Worker inside a sandboxed iframe,
// and DOM changes are communicated to the host via JSON messages.
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

func (c *RemoteDOMContent) fromWire(wire *wireUIContent) {
	c.Script = wire.Text
	c.Annotations = wire.Annotations
	// Framework is parsed from the MIME type if needed
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

func (c *BlobContent) fromWire(wire *wireUIContent) {
	if wire.Blob != "" {
		c.Data, _ = base64.StdEncoding.DecodeString(wire.Blob)
	}
	c.ContentMIMEType = wire.MIMEType
	c.Annotations = wire.Annotations
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
		c.fromWire(wire)
		return c, nil
	case wire.MIMEType == MIMETypeURLList:
		c := &URLContent{}
		c.fromWire(wire)
		return c, nil
	case len(wire.MIMEType) >= len(MIMETypeRemoteDOM) && wire.MIMEType[:len(MIMETypeRemoteDOM)] == MIMETypeRemoteDOM:
		c := &RemoteDOMContent{}
		c.fromWire(wire)
		return c, nil
	case wire.Blob != "":
		c := &BlobContent{}
		c.fromWire(wire)
		return c, nil
	default:
		return nil, fmt.Errorf("unknown content MIME type: %s", wire.MIMEType)
	}
}
