package mcp

import (
	"context"
	"strings"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ironystock/agentic-obs/internal/mcp/testutil"
)

// testPromptServer creates a minimal server for prompt testing.
func testPromptServer(t *testing.T) *Server {
	t.Helper()

	mock := testutil.NewMockOBSClient()
	mock.Connect()

	server := &Server{
		obsClient: mock,
		ctx:       context.Background(),
	}

	return server
}

// TestPromptHandlers tests that each prompt handler returns a valid result
func TestPromptHandlers(t *testing.T) {
	tests := []struct {
		name        string
		handler     func(*Server) (*mcpsdk.GetPromptResult, error)
		checkResult func(t *testing.T, result *mcpsdk.GetPromptResult)
	}{
		{
			name: "stream-launch",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleStreamLaunch(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "Pre-stream checklist")
				assert.Len(t, result.Messages, 1)
				assert.Equal(t, mcpsdk.Role("user"), result.Messages[0].Role)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "connection")
				assert.Contains(t, textContent.Text, "scenes")
				assert.Contains(t, textContent.Text, "audio")
			},
		},
		{
			name: "stream-teardown",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleStreamTeardown(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "Post-stream cleanup")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "Stop Active Streaming")
				assert.Contains(t, textContent.Text, "Offline Scene")
			},
		},
		{
			name: "audio-check",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleAudioCheck(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "audio inputs")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "Mute States")
				assert.Contains(t, textContent.Text, "Volume Levels")
			},
		},
		{
			name: "health-check",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleHealthCheck(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "comprehensive diagnostic")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "Connection Status")
				assert.Contains(t, textContent.Text, "Scene Inventory")
			},
		},
		{
			name: "recording-workflow",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleRecordingWorkflow(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "recording operations")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "Recording Status")
				assert.Contains(t, textContent.Text, "Recording Control")
			},
		},
		{
			name: "scene-organizer",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleSceneOrganizer(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "scene structure")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "List All Scenes")
				assert.Contains(t, textContent.Text, "Organization Recommendations")
			},
		},
		{
			name: "quick-status",
			handler: func(s *Server) (*mcpsdk.GetPromptResult, error) {
				return s.handleQuickStatus(context.Background(), &mcpsdk.GetPromptRequest{})
			},
			checkResult: func(t *testing.T, result *mcpsdk.GetPromptResult) {
				assert.Contains(t, result.Description, "brief summary")
				assert.Len(t, result.Messages, 1)
				textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
				assert.Contains(t, textContent.Text, "Current Scene")
				assert.Contains(t, textContent.Text, "Recording")
				assert.Contains(t, textContent.Text, "Streaming")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testPromptServer(t)
			result, err := tt.handler(server)

			require.NoError(t, err, "Handler should not return error")
			require.NotNil(t, result, "Result should not be nil")
			tt.checkResult(t, result)
		})
	}
}

// TestPromptArgumentValidation tests argument handling for prompts with required arguments
func TestPromptArgumentValidation(t *testing.T) {
	t.Run("visual-check with missing screenshot_source returns error", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{},
			},
		}

		result, err := server.handleVisualCheck(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "screenshot_source argument is required")
		assert.Nil(t, result)
	})

	t.Run("visual-check with screenshot_source succeeds", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": "main-stream",
				},
			},
		}

		result, err := server.handleVisualCheck(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Messages, 1)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Contains(t, textContent.Text, "main-stream")
	})

	t.Run("problem-detection with missing screenshot_source returns error", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{},
			},
		}

		result, err := server.handleProblemDetection(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "screenshot_source argument is required")
		assert.Nil(t, result)
	})

	t.Run("problem-detection with screenshot_source succeeds", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": "main-stream",
				},
			},
		}

		result, err := server.handleProblemDetection(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Messages, 1)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Contains(t, textContent.Text, "main-stream")
		assert.Contains(t, textContent.Text, "Black Screen Detection")
	})

	t.Run("preset-switcher without preset_name succeeds", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{},
			},
		}

		result, err := server.handlePresetSwitcher(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Messages, 1)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Contains(t, textContent.Text, "Preset Selection Guidance")
	})

	t.Run("preset-switcher with preset_name includes apply instructions", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"preset_name": "Gaming Setup",
				},
			},
		}

		result, err := server.handlePresetSwitcher(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Messages, 1)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Contains(t, textContent.Text, "Gaming Setup")
		assert.Contains(t, textContent.Text, "Apply Requested Preset")
		// Count occurrences of "Gaming Setup" - should appear twice
		count := strings.Count(textContent.Text, "Gaming Setup")
		assert.Equal(t, 2, count, "Preset name should appear twice in apply instructions")
	})
}

// TestPromptMessageContent verifies that prompt messages contain expected content
func TestPromptMessageContent(t *testing.T) {
	t.Run("stream-launch contains comprehensive checklist", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleStreamLaunch(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		// Verify key sections
		assert.Contains(t, textContent.Text, "get_obs_status")
		assert.Contains(t, textContent.Text, "list_scenes")
		assert.Contains(t, textContent.Text, "list_sources")
		assert.Contains(t, textContent.Text, "get_input_mute")
		assert.Contains(t, textContent.Text, "get_streaming_status")
		assert.Contains(t, textContent.Text, "get_recording_status")
	})

	t.Run("quick-status is concise", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleQuickStatus(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		// Verify it mentions brevity
		assert.Contains(t, textContent.Text, "brief")
		assert.Contains(t, textContent.Text, "concise")
		assert.Contains(t, textContent.Text, "3-4 lines")
	})

	t.Run("visual-check mentions screenshot analysis", func(t *testing.T) {
		server := testPromptServer(t)
		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": "test-source",
				},
			},
		}
		result, err := server.handleVisualCheck(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "Analyze Visual Layout")
		assert.Contains(t, textContent.Text, "Visual Issues")
		assert.Contains(t, textContent.Text, "test-source")
	})

	t.Run("problem-detection includes severity levels", func(t *testing.T) {
		server := testPromptServer(t)
		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": "monitor",
				},
			},
		}
		result, err := server.handleProblemDetection(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "CRITICAL")
		assert.Contains(t, textContent.Text, "WARNING")
		assert.Contains(t, textContent.Text, "Black Screen Detection")
		assert.Contains(t, textContent.Text, "Frozen Frame Detection")
	})

	t.Run("audio-check includes comprehensive checks", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleAudioCheck(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "get_input_mute")
		assert.Contains(t, textContent.Text, "get_input_volume")
		assert.Contains(t, textContent.Text, "dB")
		assert.Contains(t, textContent.Text, "multiplier")
	})

	t.Run("health-check covers all systems", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleHealthCheck(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "Connection Status")
		assert.Contains(t, textContent.Text, "Scene Inventory")
		assert.Contains(t, textContent.Text, "Source Count")
		assert.Contains(t, textContent.Text, "Recording State")
		assert.Contains(t, textContent.Text, "Streaming State")
		assert.Contains(t, textContent.Text, "Screenshot Sources")
	})

	t.Run("scene-organizer includes organization analysis", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleSceneOrganizer(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "Scene Naming Patterns")
		assert.Contains(t, textContent.Text, "Scene Redundancy")
		assert.Contains(t, textContent.Text, "Scene Completeness")
		assert.Contains(t, textContent.Text, "Organization Recommendations")
	})

	t.Run("recording-workflow includes workflow steps", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleRecordingWorkflow(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "start_recording")
		assert.Contains(t, textContent.Text, "pause_recording")
		assert.Contains(t, textContent.Text, "resume_recording")
		assert.Contains(t, textContent.Text, "stop_recording")
	})

	t.Run("stream-teardown includes cleanup steps", func(t *testing.T) {
		server := testPromptServer(t)
		result, err := server.handleStreamTeardown(context.Background(), &mcpsdk.GetPromptRequest{})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "stop_streaming")
		assert.Contains(t, textContent.Text, "stop_recording")
		assert.Contains(t, textContent.Text, "Offline")
		assert.Contains(t, textContent.Text, "Mute All Audio")
	})
}

// TestPromptMessageStructure verifies that all prompts return valid message structures
func TestPromptMessageStructure(t *testing.T) {
	server := testPromptServer(t)

	prompts := []struct {
		name    string
		handler func() (*mcpsdk.GetPromptResult, error)
	}{
		{
			name: "stream-launch",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleStreamLaunch(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "stream-teardown",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleStreamTeardown(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "audio-check",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleAudioCheck(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "health-check",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleHealthCheck(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "recording-workflow",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleRecordingWorkflow(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "scene-organizer",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleSceneOrganizer(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
		{
			name: "quick-status",
			handler: func() (*mcpsdk.GetPromptResult, error) {
				return server.handleQuickStatus(context.Background(), &mcpsdk.GetPromptRequest{})
			},
		},
	}

	for _, tt := range prompts {
		t.Run(tt.name+" has valid structure", func(t *testing.T) {
			result, err := tt.handler()

			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify description exists
			assert.NotEmpty(t, result.Description, "Description should not be empty")

			// Verify messages array
			require.NotNil(t, result.Messages, "Messages should not be nil")
			assert.Len(t, result.Messages, 1, "Should have exactly one message")

			// Verify message role
			assert.Equal(t, mcpsdk.Role("user"), result.Messages[0].Role, "Message role should be 'user'")

			// Verify message content
			require.NotNil(t, result.Messages[0].Content, "Message content should not be nil")
			textContent, ok := result.Messages[0].Content.(*mcpsdk.TextContent)
			require.True(t, ok, "Content should be TextContent type")
			assert.NotEmpty(t, textContent.Text, "Text content should not be empty")
		})
	}
}

// TestPromptArgumentInterpolation tests that arguments are properly interpolated into prompts
func TestPromptArgumentInterpolation(t *testing.T) {
	t.Run("visual-check interpolates screenshot_source", func(t *testing.T) {
		server := testPromptServer(t)

		testCases := []string{
			"main-camera",
			"stream-preview",
			"test123",
		}

		for _, sourceName := range testCases {
			req := &mcpsdk.GetPromptRequest{
				Params: &mcpsdk.GetPromptParams{
					Arguments: map[string]string{
						"screenshot_source": sourceName,
					},
				},
			}

			result, err := server.handleVisualCheck(context.Background(), req)

			require.NoError(t, err)
			textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

			// Should appear at least twice in the text
			count := strings.Count(textContent.Text, sourceName)
			assert.GreaterOrEqual(t, count, 2, "Screenshot source should appear at least twice")
		}
	})

	t.Run("problem-detection interpolates screenshot_source", func(t *testing.T) {
		server := testPromptServer(t)

		sourceName := "monitoring-source"
		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": sourceName,
				},
			},
		}

		result, err := server.handleProblemDetection(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		// Should appear at least twice in the text
		count := strings.Count(textContent.Text, sourceName)
		assert.GreaterOrEqual(t, count, 2, "Screenshot source should appear at least twice")
	})

	t.Run("preset-switcher interpolates preset_name when provided", func(t *testing.T) {
		server := testPromptServer(t)

		presetName := "StreamingSetup"
		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"preset_name": presetName,
				},
			},
		}

		result, err := server.handlePresetSwitcher(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		// Preset name should appear in the text
		assert.Contains(t, textContent.Text, presetName)
		assert.Contains(t, textContent.Text, "Apply Requested Preset")
	})

	t.Run("preset-switcher without preset_name uses default guidance", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{},
			},
		}

		result, err := server.handlePresetSwitcher(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		assert.Contains(t, textContent.Text, "Preset Selection Guidance")
		assert.NotContains(t, textContent.Text, "Apply Requested Preset")
	})
}

// TestPromptEdgeCases tests edge cases and error conditions
func TestPromptEdgeCases(t *testing.T) {
	t.Run("nil request doesn't cause panic", func(t *testing.T) {
		server := testPromptServer(t)

		// These handlers don't use req, should work with nil
		assert.NotPanics(t, func() {
			_, _ = server.handleStreamLaunch(context.Background(), nil)
		})

		assert.NotPanics(t, func() {
			_, _ = server.handleQuickStatus(context.Background(), nil)
		})
	})

	t.Run("empty screenshot_source is treated as missing", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"screenshot_source": "",
				},
			},
		}

		result, err := server.handleVisualCheck(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
		assert.Nil(t, result)
	})

	t.Run("empty preset_name is treated as not provided", func(t *testing.T) {
		server := testPromptServer(t)

		req := &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{
					"preset_name": "",
				},
			},
		}

		result, err := server.handlePresetSwitcher(context.Background(), req)

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)

		// Should use default guidance since preset_name is empty
		assert.Contains(t, textContent.Text, "Preset Selection Guidance")
	})

	t.Run("context cancellation is handled", func(t *testing.T) {
		server := testPromptServer(t)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Handlers should still work even with cancelled context
		// since they don't do async operations
		_, err := server.handleStreamLaunch(ctx, &mcpsdk.GetPromptRequest{})
		assert.NoError(t, err)
	})
}

// TestAllPromptsHaveNonEmptyMessages ensures all prompts return non-empty messages
func TestAllPromptsHaveNonEmptyMessages(t *testing.T) {
	server := testPromptServer(t)

	// Test prompts without arguments
	noArgPrompts := []struct {
		name    string
		handler func() (*mcpsdk.GetPromptResult, error)
	}{
		{"stream-launch", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleStreamLaunch(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"stream-teardown", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleStreamTeardown(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"audio-check", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleAudioCheck(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"health-check", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleHealthCheck(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"recording-workflow", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleRecordingWorkflow(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"scene-organizer", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleSceneOrganizer(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"quick-status", func() (*mcpsdk.GetPromptResult, error) {
			return server.handleQuickStatus(context.Background(), &mcpsdk.GetPromptRequest{})
		}},
		{"preset-switcher", func() (*mcpsdk.GetPromptResult, error) {
			return server.handlePresetSwitcher(context.Background(), &mcpsdk.GetPromptRequest{
				Params: &mcpsdk.GetPromptParams{Arguments: map[string]string{}},
			})
		}},
	}

	for _, tt := range noArgPrompts {
		t.Run(tt.name+" has non-empty message", func(t *testing.T) {
			result, err := tt.handler()

			require.NoError(t, err)
			require.NotNil(t, result)
			require.Len(t, result.Messages, 1)

			textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
			assert.NotEmpty(t, textContent.Text)
			// All messages should be at least 100 characters (reasonable minimum)
			assert.Greater(t, len(textContent.Text), 100)
		})
	}

	// Test prompts with required arguments
	t.Run("visual-check has non-empty message", func(t *testing.T) {
		result, err := server.handleVisualCheck(context.Background(), &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{"screenshot_source": "test"},
			},
		})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Greater(t, len(textContent.Text), 100)
	})

	t.Run("problem-detection has non-empty message", func(t *testing.T) {
		result, err := server.handleProblemDetection(context.Background(), &mcpsdk.GetPromptRequest{
			Params: &mcpsdk.GetPromptParams{
				Arguments: map[string]string{"screenshot_source": "test"},
			},
		})

		require.NoError(t, err)
		textContent := result.Messages[0].Content.(*mcpsdk.TextContent)
		assert.Greater(t, len(textContent.Text), 100)
	})
}
