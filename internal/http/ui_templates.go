package http

import (
	"embed"
	"html/template"
	"sync"
)

//go:embed templates/*.html templates/*.css
var templateFiles embed.FS

// templateCache holds parsed templates
var (
	templateCache     map[string]*template.Template
	templateCacheMu   sync.RWMutex
	templateCacheOnce sync.Once
)

// initTemplateCache initializes the template cache
func initTemplateCache() {
	templateCache = make(map[string]*template.Template)
}

// getSharedCSS reads the shared CSS file
func getSharedCSS() string {
	data, err := templateFiles.ReadFile("templates/shared.css")
	if err != nil {
		return ""
	}
	return string(data)
}

// getTemplate returns a parsed template by name, with caching
func getTemplate(name string) (*template.Template, error) {
	templateCacheOnce.Do(initTemplateCache)

	templateCacheMu.RLock()
	if tmpl, ok := templateCache[name]; ok {
		templateCacheMu.RUnlock()
		return tmpl, nil
	}
	templateCacheMu.RUnlock()

	// Parse template
	data, err := templateFiles.ReadFile("templates/" + name + ".html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}

	// Cache the template
	templateCacheMu.Lock()
	templateCache[name] = tmpl
	templateCacheMu.Unlock()

	return tmpl, nil
}

// Template names
const (
	templateStatusDashboard  = "status_dashboard"
	templateScenePreview     = "scene_preview"
	templateAudioMixer       = "audio_mixer"
	templateScreenshotGallery = "screenshot_gallery"
	templateError            = "error"
)
