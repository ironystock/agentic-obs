package mcp

import (
	"fmt"

	"github.com/ironystock/mcpui-go"
)

// UI resource URI patterns
const (
	UIStatusDashboardURI   = "ui://status/dashboard"
	UIScenePreviewURI      = "ui://scene/preview"
	UIAudioMixerURI        = "ui://audio/mixer"
	UIScreenshotGalleryURI = "ui://screenshot/gallery"
)

// UIResourceRegistry defines available UI resources for MCP clients.
type UIResourceRegistry struct {
	httpBaseURL string
}

// NewUIResourceRegistry creates a new registry with the HTTP server base URL.
func NewUIResourceRegistry(httpBaseURL string) *UIResourceRegistry {
	return &UIResourceRegistry{
		httpBaseURL: httpBaseURL,
	}
}

// ListUIResources returns all available UI resources.
func (r *UIResourceRegistry) ListUIResources() []*mcpui.UIResource {
	return []*mcpui.UIResource{
		{
			URI:         UIStatusDashboardURI,
			Name:        "status-dashboard",
			Title:       "OBS Status Dashboard",
			Description: "Real-time OBS connection status, recording/streaming state, and scene overview",
			MIMEType:    mcpui.MIMETypeHTML,
		},
		{
			URI:         UIScenePreviewURI,
			Name:        "scene-preview",
			Title:       "Scene Preview",
			Description: "Visual scene grid with thumbnails for quick scene switching",
			MIMEType:    mcpui.MIMETypeHTML,
		},
		{
			URI:         UIAudioMixerURI,
			Name:        "audio-mixer",
			Title:       "Audio Mixer",
			Description: "Audio input controls with volume sliders and mute toggles",
			MIMEType:    mcpui.MIMETypeHTML,
		},
		{
			URI:         UIScreenshotGalleryURI,
			Name:        "screenshot-gallery",
			Title:       "Screenshot Gallery",
			Description: "Live screenshot gallery from configured screenshot sources",
			MIMEType:    mcpui.MIMETypeHTML,
		},
	}
}

// GetUIResourceTemplate returns a template for dynamic UI resources.
func (r *UIResourceRegistry) GetUIResourceTemplate(uri string) *mcpui.UIResourceTemplate {
	switch uri {
	case UIStatusDashboardURI:
		return &mcpui.UIResourceTemplate{
			URITemplate: UIStatusDashboardURI,
			Name:        "status-dashboard",
			Title:       "OBS Status Dashboard",
			Description: "Real-time OBS status information",
			MIMEType:    mcpui.MIMETypeHTML,
		}
	default:
		return nil
	}
}

// BuildStatusDashboardURL returns the HTTP URL for the status dashboard UI.
func (r *UIResourceRegistry) BuildStatusDashboardURL() string {
	return fmt.Sprintf("%s/ui/status", r.httpBaseURL)
}

// BuildScenePreviewURL returns the HTTP URL for the scene preview UI.
func (r *UIResourceRegistry) BuildScenePreviewURL() string {
	return fmt.Sprintf("%s/ui/scenes", r.httpBaseURL)
}

// BuildAudioMixerURL returns the HTTP URL for the audio mixer UI.
func (r *UIResourceRegistry) BuildAudioMixerURL() string {
	return fmt.Sprintf("%s/ui/audio", r.httpBaseURL)
}

// BuildScreenshotGalleryURL returns the HTTP URL for the screenshot gallery UI.
func (r *UIResourceRegistry) BuildScreenshotGalleryURL() string {
	return fmt.Sprintf("%s/ui/screenshots", r.httpBaseURL)
}

// NewStatusDashboardResource creates a UIResourceContents for the status dashboard.
// This returns an URLContent pointing to the HTTP-served UI.
func (r *UIResourceRegistry) NewStatusDashboardResource() (*mcpui.UIResourceContents, error) {
	content := &mcpui.URLContent{
		URL: r.BuildStatusDashboardURL(),
	}
	return mcpui.NewUIResourceContents(UIStatusDashboardURI, content)
}

// NewScenePreviewResource creates a UIResourceContents for the scene preview.
func (r *UIResourceRegistry) NewScenePreviewResource() (*mcpui.UIResourceContents, error) {
	content := &mcpui.URLContent{
		URL: r.BuildScenePreviewURL(),
	}
	return mcpui.NewUIResourceContents(UIScenePreviewURI, content)
}

// NewAudioMixerResource creates a UIResourceContents for the audio mixer.
func (r *UIResourceRegistry) NewAudioMixerResource() (*mcpui.UIResourceContents, error) {
	content := &mcpui.URLContent{
		URL: r.BuildAudioMixerURL(),
	}
	return mcpui.NewUIResourceContents(UIAudioMixerURI, content)
}

// NewScreenshotGalleryResource creates a UIResourceContents for the screenshot gallery.
func (r *UIResourceRegistry) NewScreenshotGalleryResource() (*mcpui.UIResourceContents, error) {
	content := &mcpui.URLContent{
		URL: r.BuildScreenshotGalleryURL(),
	}
	return mcpui.NewUIResourceContents(UIScreenshotGalleryURI, content)
}
