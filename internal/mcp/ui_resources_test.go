package mcp

import (
	"testing"

	"github.com/ironystock/agentic-obs/pkg/mcpui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUIResourceRegistry(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	assert.NotNil(t, registry)
	assert.Equal(t, "http://localhost:8765", registry.httpBaseURL)
}

func TestUIResourceRegistry_ListUIResources(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	resources := registry.ListUIResources()

	assert.Len(t, resources, 4)

	// Verify all expected resources are present
	uriMap := make(map[string]*mcpui.UIResource)
	for _, r := range resources {
		uriMap[r.URI] = r
	}

	// Status dashboard
	assert.Contains(t, uriMap, UIStatusDashboardURI)
	assert.Equal(t, "status-dashboard", uriMap[UIStatusDashboardURI].Name)
	assert.Equal(t, mcpui.MIMETypeHTML, uriMap[UIStatusDashboardURI].MIMEType)

	// Scene preview
	assert.Contains(t, uriMap, UIScenePreviewURI)
	assert.Equal(t, "scene-preview", uriMap[UIScenePreviewURI].Name)

	// Audio mixer
	assert.Contains(t, uriMap, UIAudioMixerURI)
	assert.Equal(t, "audio-mixer", uriMap[UIAudioMixerURI].Name)

	// Screenshot gallery
	assert.Contains(t, uriMap, UIScreenshotGalleryURI)
	assert.Equal(t, "screenshot-gallery", uriMap[UIScreenshotGalleryURI].Name)
}

func TestUIResourceRegistry_BuildURLs(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:9000")

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "status dashboard URL",
			method:   registry.BuildStatusDashboardURL,
			expected: "http://localhost:9000/ui/status",
		},
		{
			name:     "scene preview URL",
			method:   registry.BuildScenePreviewURL,
			expected: "http://localhost:9000/ui/scenes",
		},
		{
			name:     "audio mixer URL",
			method:   registry.BuildAudioMixerURL,
			expected: "http://localhost:9000/ui/audio",
		},
		{
			name:     "screenshot gallery URL",
			method:   registry.BuildScreenshotGalleryURL,
			expected: "http://localhost:9000/ui/screenshots",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method())
		})
	}
}

func TestUIResourceRegistry_NewStatusDashboardResource(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	resource, err := registry.NewStatusDashboardResource()

	require.NoError(t, err)
	assert.Equal(t, UIStatusDashboardURI, resource.URI)
	assert.Equal(t, mcpui.MIMETypeURLList, resource.MIMEType)
	assert.Equal(t, "http://localhost:8765/ui/status", resource.Text)
}

func TestUIResourceRegistry_NewScenePreviewResource(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	resource, err := registry.NewScenePreviewResource()

	require.NoError(t, err)
	assert.Equal(t, UIScenePreviewURI, resource.URI)
	assert.Equal(t, mcpui.MIMETypeURLList, resource.MIMEType)
	assert.Equal(t, "http://localhost:8765/ui/scenes", resource.Text)
}

func TestUIResourceRegistry_NewAudioMixerResource(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	resource, err := registry.NewAudioMixerResource()

	require.NoError(t, err)
	assert.Equal(t, UIAudioMixerURI, resource.URI)
	assert.Equal(t, mcpui.MIMETypeURLList, resource.MIMEType)
	assert.Equal(t, "http://localhost:8765/ui/audio", resource.Text)
}

func TestUIResourceRegistry_NewScreenshotGalleryResource(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	resource, err := registry.NewScreenshotGalleryResource()

	require.NoError(t, err)
	assert.Equal(t, UIScreenshotGalleryURI, resource.URI)
	assert.Equal(t, mcpui.MIMETypeURLList, resource.MIMEType)
	assert.Equal(t, "http://localhost:8765/ui/screenshots", resource.Text)
}

func TestUIResourceRegistry_GetUIResourceTemplate(t *testing.T) {
	registry := NewUIResourceRegistry("http://localhost:8765")

	t.Run("returns template for status dashboard", func(t *testing.T) {
		template := registry.GetUIResourceTemplate(UIStatusDashboardURI)

		require.NotNil(t, template)
		assert.Equal(t, UIStatusDashboardURI, template.URITemplate)
		assert.Equal(t, "status-dashboard", template.Name)
		assert.Equal(t, mcpui.MIMETypeHTML, template.MIMEType)
	})

	t.Run("returns nil for unknown URI", func(t *testing.T) {
		template := registry.GetUIResourceTemplate("ui://unknown/resource")

		assert.Nil(t, template)
	})
}

func TestUIResourceURIConstants(t *testing.T) {
	// Verify URI constants follow the ui:// scheme
	assert.True(t, len(UIStatusDashboardURI) > 5)
	assert.Equal(t, "ui://", UIStatusDashboardURI[:5])

	assert.True(t, len(UIScenePreviewURI) > 5)
	assert.Equal(t, "ui://", UIScenePreviewURI[:5])

	assert.True(t, len(UIAudioMixerURI) > 5)
	assert.Equal(t, "ui://", UIAudioMixerURI[:5])

	assert.True(t, len(UIScreenshotGalleryURI) > 5)
	assert.Equal(t, "ui://", UIScreenshotGalleryURI[:5])
}
