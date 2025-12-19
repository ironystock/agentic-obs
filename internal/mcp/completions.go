package mcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

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
// Supports preset_name and screenshot_source arguments across relevant prompts.
func (s *Server) completePromptArg(ctx context.Context, promptName, argName, value string) (*mcpsdk.CompleteResult, error) {
	log.Printf("Completing prompt argument: prompt=%s, arg=%s, value=%s", promptName, argName, value)

	var completions []string
	var err error

	// Handle preset_name argument (used by preset-switcher prompt)
	if argName == "preset_name" {
		completions, err = s.getPresetNameCompletions(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get preset completions: %w", err)
		}
	}

	// Handle screenshot_source argument (used by visual-check and problem-detection prompts)
	if argName == "screenshot_source" {
		completions, err = s.getScreenshotSourceCompletions(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get screenshot source completions: %w", err)
		}
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
// Supports obs://scene/, obs://preset/, and obs://screenshot/ URIs.
func (s *Server) completeResourceURI(ctx context.Context, uri, value string) (*mcpsdk.CompleteResult, error) {
	log.Printf("Completing resource URI: uri=%s, value=%s", uri, value)

	var completions []string
	var err error

	// Determine resource type from URI prefix
	switch {
	case strings.HasPrefix(uri, SceneURIPrefix):
		// Complete scene names
		completions, err = s.getSceneNameCompletions(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get scene completions: %w", err)
		}

	case strings.HasPrefix(uri, PresetURIPrefix):
		// Complete preset names
		completions, err = s.getPresetNameCompletions(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get preset completions: %w", err)
		}

	case strings.HasPrefix(uri, ScreenshotURIPrefix):
		// Complete screenshot source names
		completions, err = s.getScreenshotSourceCompletions(ctx, value)
		if err != nil {
			return nil, fmt.Errorf("failed to get screenshot source completions: %w", err)
		}

	case strings.HasPrefix(uri, ScreenshotURLURIPrefix):
		// Complete screenshot source names (for URL resources)
		completions, err = s.getScreenshotSourceCompletions(ctx, value)
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

// getSceneNameCompletions fetches all scene names from OBS.
func (s *Server) getSceneNameCompletions(ctx context.Context, prefix string) ([]string, error) {
	scenes, _, err := s.obsClient.GetSceneList()
	if err != nil {
		return nil, fmt.Errorf("failed to get scene list: %w", err)
	}
	return scenes, nil
}

// getPresetNameCompletions fetches all preset names from storage.
func (s *Server) getPresetNameCompletions(ctx context.Context, prefix string) ([]string, error) {
	presets, err := s.storage.ListScenePresets(ctx, "") // Empty filter gets all presets
	if err != nil {
		return nil, fmt.Errorf("failed to list presets: %w", err)
	}

	names := make([]string, len(presets))
	for i, preset := range presets {
		names[i] = preset.Name
	}
	return names, nil
}

// getScreenshotSourceCompletions fetches all screenshot source names from storage.
func (s *Server) getScreenshotSourceCompletions(ctx context.Context, prefix string) ([]string, error) {
	sources, err := s.storage.ListScreenshotSources(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list screenshot sources: %w", err)
	}

	names := make([]string, len(sources))
	for i, source := range sources {
		names[i] = source.Name
	}
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
