package mcp

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// completionCache provides short-lived caching for completion results.
// This reduces repeated calls to OBS/storage for autocomplete scenarios
// where users type multiple characters in quick succession.
type completionCache struct {
	mu            sync.RWMutex
	scenes        []string
	presets       []string
	sources       []string
	scenesTTL     time.Time
	presetsTTL    time.Time
	sourcesTTL    time.Time
	cacheDuration time.Duration
}

// newCompletionCache creates a cache with the specified TTL.
func newCompletionCache(ttl time.Duration) *completionCache {
	return &completionCache{
		cacheDuration: ttl,
	}
}

// Global completion cache with 5-second TTL
var compCache = newCompletionCache(5 * time.Second)

// resetCompletionCache clears the cache (used in tests)
func resetCompletionCache() {
	compCache.mu.Lock()
	defer compCache.mu.Unlock()
	compCache.scenes = nil
	compCache.presets = nil
	compCache.sources = nil
	compCache.scenesTTL = time.Time{}
	compCache.presetsTTL = time.Time{}
	compCache.sourcesTTL = time.Time{}
}

// handleCompletion dispatches completion requests to the appropriate handler
// based on the reference type (prompt argument or resource URI).
func (s *Server) handleCompletion(ctx context.Context, req *mcpsdk.CompleteRequest) (*mcpsdk.CompleteResult, error) {
	if req == nil || req.Params == nil || req.Params.Ref == nil {
		return nil, fmt.Errorf("invalid completion request: missing ref")
	}

	refType := req.Params.Ref.Type
	argName := req.Params.Argument.Name
	value := req.Params.Argument.Value

	log.Printf("Handling completion: ref_type=%s, arg_name=%s, value=%s", refType, argName, value)

	switch refType {
	case "ref/prompt":
		// Complete prompt arguments
		promptName := req.Params.Ref.Name
		if promptName == "" {
			return nil, fmt.Errorf("prompt name required for ref/prompt completion")
		}
		return s.completePromptArg(ctx, promptName, argName, value)

	case "ref/resource":
		// Complete resource URIs
		uri := req.Params.Ref.URI
		if uri == "" {
			return nil, fmt.Errorf("URI required for ref/resource completion")
		}
		return s.completeResourceURI(ctx, uri, value)

	default:
		return nil, fmt.Errorf("unsupported completion reference type: %s", refType)
	}
}

// completePromptArg handles completions for prompt arguments.
// Supports the following arguments:
//   - preset_name: Used by preset-switcher prompt
//   - screenshot_source: Used by visual-check and problem-detection prompts
//   - scene_name: Used by scene-designer and source-management prompts
//   - monitor_scene: Used by visual-setup prompt
func (s *Server) completePromptArg(ctx context.Context, promptName, argName, value string) (*mcpsdk.CompleteResult, error) {
	log.Printf("Completing prompt argument: prompt=%s, arg=%s, value=%s", promptName, argName, value)

	var completions []string
	var err error

	switch argName {
	case "preset_name":
		// Used by preset-switcher prompt
		completions, err = s.getPresetNameCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get preset completions: %w", err)
		}

	case "screenshot_source":
		// Used by visual-check and problem-detection prompts
		completions, err = s.getScreenshotSourceCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get screenshot source completions: %w", err)
		}

	case "scene_name", "monitor_scene":
		// scene_name: Used by scene-designer and source-management prompts
		// monitor_scene: Used by visual-setup prompt
		completions, err = s.getSceneNameCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get scene completions: %w", err)
		}

	default:
		// Unknown argument - return empty completions
		log.Printf("No completions available for argument: %s", argName)
		completions = []string{}
	}

	// Filter completions by prefix
	filtered := filterCompletions(completions, value)

	return &mcpsdk.CompleteResult{
		Completion: mcpsdk.CompletionResultDetails{
			Values:  filtered,
			Total:   len(filtered),
			HasMore: false,
		},
	}, nil
}

// completeResourceURI handles completions for resource URIs.
// Supports obs://scene/, obs://preset/, obs://screenshot/, and obs://screenshot-url/ URIs.
func (s *Server) completeResourceURI(ctx context.Context, uri, value string) (*mcpsdk.CompleteResult, error) {
	log.Printf("Completing resource URI: uri=%s, value=%s", uri, value)

	var completions []string
	var err error

	// Determine resource type from URI prefix
	switch {
	case strings.HasPrefix(uri, SceneURIPrefix):
		// Complete scene names
		completions, err = s.getSceneNameCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get scene completions: %w", err)
		}

	case strings.HasPrefix(uri, PresetURIPrefix):
		// Complete preset names
		completions, err = s.getPresetNameCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get preset completions: %w", err)
		}

	case strings.HasPrefix(uri, ScreenshotURIPrefix), strings.HasPrefix(uri, ScreenshotURLURIPrefix):
		// Complete screenshot source names (both binary and URL resources)
		completions, err = s.getScreenshotSourceCompletions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get screenshot source completions: %w", err)
		}

	default:
		// Unknown resource type - return empty completions
		log.Printf("No completions available for URI: %s", uri)
		return &mcpsdk.CompleteResult{
			Completion: mcpsdk.CompletionResultDetails{
				Values:  []string{},
				Total:   0,
				HasMore: false,
			},
		}, nil
	}

	// Filter completions by prefix
	filtered := filterCompletions(completions, value)

	return &mcpsdk.CompleteResult{
		Completion: mcpsdk.CompletionResultDetails{
			Values:  filtered,
			Total:   len(filtered),
			HasMore: false,
		},
	}, nil
}

// getSceneNameCompletions fetches all scene names from OBS with caching.
// Results are cached for 5 seconds to reduce OBS API calls during rapid typing.
func (s *Server) getSceneNameCompletions(ctx context.Context) ([]string, error) {
	compCache.mu.RLock()
	if time.Now().Before(compCache.scenesTTL) && compCache.scenes != nil {
		scenes := compCache.scenes
		compCache.mu.RUnlock()
		log.Printf("Using cached scene completions (%d scenes)", len(scenes))
		return scenes, nil
	}
	compCache.mu.RUnlock()

	// Cache miss - fetch from OBS
	scenes, _, err := s.obsClient.GetSceneList()
	if err != nil {
		return nil, fmt.Errorf("failed to get scene list: %w", err)
	}

	// Update cache
	compCache.mu.Lock()
	compCache.scenes = scenes
	compCache.scenesTTL = time.Now().Add(compCache.cacheDuration)
	compCache.mu.Unlock()

	log.Printf("Fetched and cached scene completions (%d scenes)", len(scenes))
	return scenes, nil
}

// getPresetNameCompletions fetches all preset names from storage with caching.
// Results are cached for 5 seconds to reduce storage calls during rapid typing.
func (s *Server) getPresetNameCompletions(ctx context.Context) ([]string, error) {
	compCache.mu.RLock()
	if time.Now().Before(compCache.presetsTTL) && compCache.presets != nil {
		presets := compCache.presets
		compCache.mu.RUnlock()
		log.Printf("Using cached preset completions (%d presets)", len(presets))
		return presets, nil
	}
	compCache.mu.RUnlock()

	// Cache miss - fetch from storage
	presets, err := s.storage.ListScenePresets(ctx, "") // Empty filter gets all presets
	if err != nil {
		return nil, fmt.Errorf("failed to list presets: %w", err)
	}

	names := make([]string, len(presets))
	for i, preset := range presets {
		names[i] = preset.Name
	}

	// Update cache
	compCache.mu.Lock()
	compCache.presets = names
	compCache.presetsTTL = time.Now().Add(compCache.cacheDuration)
	compCache.mu.Unlock()

	log.Printf("Fetched and cached preset completions (%d presets)", len(names))
	return names, nil
}

// getScreenshotSourceCompletions fetches all screenshot source names from storage with caching.
// Results are cached for 5 seconds to reduce storage calls during rapid typing.
func (s *Server) getScreenshotSourceCompletions(ctx context.Context) ([]string, error) {
	compCache.mu.RLock()
	if time.Now().Before(compCache.sourcesTTL) && compCache.sources != nil {
		sources := compCache.sources
		compCache.mu.RUnlock()
		log.Printf("Using cached screenshot source completions (%d sources)", len(sources))
		return sources, nil
	}
	compCache.mu.RUnlock()

	// Cache miss - fetch from storage
	sources, err := s.storage.ListScreenshotSources(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list screenshot sources: %w", err)
	}

	names := make([]string, len(sources))
	for i, source := range sources {
		names[i] = source.Name
	}

	// Update cache
	compCache.mu.Lock()
	compCache.sources = names
	compCache.sourcesTTL = time.Now().Add(compCache.cacheDuration)
	compCache.mu.Unlock()

	log.Printf("Fetched and cached screenshot source completions (%d sources)", len(names))
	return names, nil
}

// filterCompletions filters a list of completions by prefix (case-insensitive).
// Returns all completions if prefix is empty.
func filterCompletions(completions []string, prefix string) []string {
	if prefix == "" {
		return completions
	}

	lowerPrefix := strings.ToLower(prefix)
	filtered := make([]string, 0, len(completions))

	for _, completion := range completions {
		if strings.HasPrefix(strings.ToLower(completion), lowerPrefix) {
			filtered = append(filtered, completion)
		}
	}

	return filtered
}
