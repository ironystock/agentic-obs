package mcp

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ironystock/agentic-obs/internal/mcp/testutil"
	"github.com/ironystock/agentic-obs/internal/storage"
)

// testServerForToolConfig creates a test server with mock OBS client and storage for tool config tests.
func testServerForToolConfig(t *testing.T) (*Server, *testutil.MockOBSClient, *storage.DB) {
	t.Helper()

	mock := testutil.NewMockOBSClient()
	mock.Connect()

	// Create temp database
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := storage.New(context.Background(), storage.Config{Path: dbPath})
	require.NoError(t, err)

	server := &Server{
		obsClient:  mock,
		storage:    db,
		toolGroups: DefaultToolGroupConfig(), // All groups enabled
		ctx:        context.Background(),
	}

	t.Cleanup(func() {
		db.Close()
	})

	return server, mock, db
}

// Test get_tool_config handler

func TestHandleGetToolConfig(t *testing.T) {
	t.Run("returns all groups when no filter", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := GetToolConfigInput{}
		_, result, err := server.handleGetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok, "result should be a map")

		groups, ok := resultMap["groups"].([]ToolGroupInfo)
		require.True(t, ok, "groups should be []ToolGroupInfo")
		assert.Len(t, groups, 8, "should have 8 tool groups")

		// Verify all groups are enabled by default
		for _, g := range groups {
			assert.True(t, g.Enabled, "group %s should be enabled", g.Name)
		}
	})

	t.Run("filters by group name", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := GetToolConfigInput{Group: "Audio"}
		_, result, err := server.handleGetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 1, "should have 1 group when filtering")
		assert.Equal(t, "Audio", groups[0].Name)
		assert.Equal(t, 4, groups[0].ToolCount)
	})

	t.Run("includes tool names when verbose", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := GetToolConfigInput{Group: "Audio", Verbose: true}
		_, result, err := server.handleGetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 1)
		assert.NotNil(t, groups[0].Tools, "should include tool names")
		assert.Contains(t, groups[0].Tools, "toggle_input_mute")
	})

	t.Run("excludes tool names when not verbose", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := GetToolConfigInput{Group: "Audio", Verbose: false}
		_, result, err := server.handleGetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 1)
		assert.Nil(t, groups[0].Tools, "should not include tool names")
	})

	t.Run("includes meta tools in count", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := GetToolConfigInput{}
		_, result, err := server.handleGetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})

		metaTools := resultMap["meta_tools"].([]string)
		assert.Len(t, metaTools, 4, "should have 4 meta tools")
		assert.Contains(t, metaTools, "help")
		assert.Contains(t, metaTools, "get_tool_config")
		assert.Contains(t, metaTools, "set_tool_config")
		assert.Contains(t, metaTools, "list_tool_groups")
	})
}

// Test set_tool_config handler

func TestHandleSetToolConfig(t *testing.T) {
	t.Run("disables tool group", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		// Verify Audio is enabled initially
		assert.True(t, server.toolGroups.Audio)

		input := SetToolConfigInput{Group: "Audio", Enabled: false}
		_, result, err := server.handleSetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap := result.(map[string]interface{})
		assert.Equal(t, "Audio", resultMap["group"])
		assert.Equal(t, true, resultMap["previous_state"])
		assert.Equal(t, false, resultMap["new_state"])
		assert.Equal(t, 4, resultMap["tools_affected"])

		// Verify group is now disabled
		assert.False(t, server.toolGroups.Audio)
	})

	t.Run("enables tool group", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		// First disable the group
		server.toolGroups.Visual = false

		input := SetToolConfigInput{Group: "Visual", Enabled: true}
		_, result, err := server.handleSetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})

		assert.Equal(t, false, resultMap["previous_state"])
		assert.Equal(t, true, resultMap["new_state"])

		// Verify group is now enabled
		assert.True(t, server.toolGroups.Visual)
	})

	t.Run("rejects invalid group name", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := SetToolConfigInput{Group: "InvalidGroup", Enabled: false}
		_, _, err := server.handleSetToolConfig(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid group name")
		assert.Contains(t, err.Error(), "InvalidGroup")
	})

	t.Run("persists to database when requested", func(t *testing.T) {
		server, _, db := testServerForToolConfig(t)

		input := SetToolConfigInput{Group: "Layout", Enabled: false, Persist: true}
		_, result, err := server.handleSetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, true, resultMap["persisted"])

		// Verify config was persisted
		config, err := db.LoadToolGroupConfig(context.Background())
		require.NoError(t, err)
		assert.False(t, config.Layout, "Layout should be persisted as disabled")
	})

	t.Run("does not persist by default", func(t *testing.T) {
		server, _, db := testServerForToolConfig(t)

		// Set initial state in DB
		err := db.SaveToolGroupConfig(context.Background(), storage.ToolGroupConfig{
			Core: true, Visual: true, Layout: true, Audio: true,
			Sources: true, Design: true, Filters: true, Transitions: true,
		})
		require.NoError(t, err)

		input := SetToolConfigInput{Group: "Layout", Enabled: false, Persist: false}
		_, result, err := server.handleSetToolConfig(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, false, resultMap["persisted"])

		// Verify config was NOT changed in DB
		config, err := db.LoadToolGroupConfig(context.Background())
		require.NoError(t, err)
		assert.True(t, config.Layout, "Layout should still be enabled in DB")
	})

	t.Run("handles all tool groups", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		for _, group := range ToolGroupOrder {
			input := SetToolConfigInput{Group: group, Enabled: false}
			_, _, err := server.handleSetToolConfig(context.Background(), nil, input)
			assert.NoError(t, err, "should handle group %s", group)
		}

		// Verify all groups are disabled
		assert.False(t, server.toolGroups.Core)
		assert.False(t, server.toolGroups.Sources)
		assert.False(t, server.toolGroups.Audio)
		assert.False(t, server.toolGroups.Layout)
		assert.False(t, server.toolGroups.Visual)
		assert.False(t, server.toolGroups.Design)
		assert.False(t, server.toolGroups.Filters)
		assert.False(t, server.toolGroups.Transitions)
	})
}

// Test list_tool_groups handler

func TestHandleListToolGroups(t *testing.T) {
	t.Run("lists all groups by default", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := ListToolGroupsInput{}
		_, result, err := server.handleListToolGroups(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 8, "should list all 8 groups")
		assert.Equal(t, 8, resultMap["count"])

		// Verify correct order
		expectedOrder := []string{"Core", "Sources", "Audio", "Layout", "Visual", "Design", "Filters", "Transitions"}
		for i, expectedName := range expectedOrder {
			assert.Equal(t, expectedName, groups[i].Name, "group %d should be %s", i, expectedName)
		}
	})

	t.Run("includes disabled groups by default", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		// Disable some groups
		server.toolGroups.Audio = false
		server.toolGroups.Visual = false

		input := ListToolGroupsInput{IncludeDisabled: true}
		_, result, err := server.handleListToolGroups(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 8, "should include disabled groups")

		// Verify Audio and Visual show as disabled
		var audioFound, visualFound bool
		for _, g := range groups {
			if g.Name == "Audio" {
				audioFound = true
				assert.False(t, g.Enabled)
			}
			if g.Name == "Visual" {
				visualFound = true
				assert.False(t, g.Enabled)
			}
		}
		assert.True(t, audioFound, "Audio group should be listed")
		assert.True(t, visualFound, "Visual group should be listed")
	})

	t.Run("excludes disabled groups when requested", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		// Disable some groups
		server.toolGroups.Audio = false
		server.toolGroups.Visual = false

		input := ListToolGroupsInput{IncludeDisabled: false}
		_, result, err := server.handleListToolGroups(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		groups := resultMap["groups"].([]ToolGroupInfo)

		assert.Len(t, groups, 6, "should exclude 2 disabled groups")

		// Verify Audio and Visual are not in the list
		for _, g := range groups {
			assert.NotEqual(t, "Audio", g.Name, "Audio should not be in list")
			assert.NotEqual(t, "Visual", g.Name, "Visual should not be in list")
		}
	})

	t.Run("includes meta tools in response", func(t *testing.T) {
		server, _, _ := testServerForToolConfig(t)

		input := ListToolGroupsInput{}
		_, result, err := server.handleListToolGroups(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})

		metaTools := resultMap["meta_tools"].([]string)
		assert.Len(t, metaTools, 4)
	})
}

// Test getGroupEnabled and setGroupEnabled helpers

func TestGetGroupEnabled(t *testing.T) {
	server, _, _ := testServerForToolConfig(t)

	testCases := []struct {
		group    string
		field    *bool
		expected bool
	}{
		{"Core", &server.toolGroups.Core, true},
		{"Sources", &server.toolGroups.Sources, true},
		{"Audio", &server.toolGroups.Audio, true},
		{"Layout", &server.toolGroups.Layout, true},
		{"Visual", &server.toolGroups.Visual, true},
		{"Design", &server.toolGroups.Design, true},
		{"Filters", &server.toolGroups.Filters, true},
		{"Transitions", &server.toolGroups.Transitions, true},
		{"Unknown", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.group, func(t *testing.T) {
			result := server.getGroupEnabled(tc.group)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSetGroupEnabled(t *testing.T) {
	server, _, _ := testServerForToolConfig(t)

	// Disable all groups
	for _, group := range ToolGroupOrder {
		server.setGroupEnabled(group, false)
	}

	// Verify all disabled
	assert.False(t, server.toolGroups.Core)
	assert.False(t, server.toolGroups.Sources)
	assert.False(t, server.toolGroups.Audio)
	assert.False(t, server.toolGroups.Layout)
	assert.False(t, server.toolGroups.Visual)
	assert.False(t, server.toolGroups.Design)
	assert.False(t, server.toolGroups.Filters)
	assert.False(t, server.toolGroups.Transitions)

	// Re-enable all
	for _, group := range ToolGroupOrder {
		server.setGroupEnabled(group, true)
	}

	// Verify all enabled
	assert.True(t, server.toolGroups.Core)
	assert.True(t, server.toolGroups.Sources)
	assert.True(t, server.toolGroups.Audio)
	assert.True(t, server.toolGroups.Layout)
	assert.True(t, server.toolGroups.Visual)
	assert.True(t, server.toolGroups.Design)
	assert.True(t, server.toolGroups.Filters)
	assert.True(t, server.toolGroups.Transitions)
}

func TestConvertToStorageConfig(t *testing.T) {
	server, _, _ := testServerForToolConfig(t)

	// Set some groups disabled
	server.toolGroups.Audio = false
	server.toolGroups.Filters = false

	config := server.convertToStorageConfig()

	assert.True(t, config.Core)
	assert.True(t, config.Sources)
	assert.False(t, config.Audio)
	assert.True(t, config.Layout)
	assert.True(t, config.Visual)
	assert.True(t, config.Design)
	assert.False(t, config.Filters)
	assert.True(t, config.Transitions)
}

// Test tool group metadata

func TestToolGroupMetadata(t *testing.T) {
	expectedGroups := map[string]struct {
		toolCount int
		hasTools  []string
	}{
		"Core": {
			toolCount: 25,
			hasTools:  []string{"list_scenes", "start_recording", "toggle_virtual_cam", "toggle_studio_mode"},
		},
		"Sources": {
			toolCount: 3,
			hasTools:  []string{"list_sources", "toggle_source_visibility"},
		},
		"Audio": {
			toolCount: 4,
			hasTools:  []string{"toggle_input_mute", "set_input_volume"},
		},
		"Layout": {
			toolCount: 6,
			hasTools:  []string{"save_scene_preset", "apply_scene_preset"},
		},
		"Visual": {
			toolCount: 4,
			hasTools:  []string{"create_screenshot_source", "list_screenshot_sources"},
		},
		"Design": {
			toolCount: 14,
			hasTools:  []string{"create_text_source", "set_source_transform"},
		},
		"Filters": {
			toolCount: 7,
			hasTools:  []string{"list_source_filters", "toggle_source_filter"},
		},
		"Transitions": {
			toolCount: 5,
			hasTools:  []string{"list_transitions", "set_current_transition"},
		},
	}

	for groupName, expected := range expectedGroups {
		t.Run(groupName, func(t *testing.T) {
			meta := toolGroupMetadata[groupName]
			require.NotNil(t, meta, "metadata should exist for %s", groupName)

			assert.Equal(t, expected.toolCount, meta.ToolCount, "tool count mismatch")
			assert.Len(t, meta.ToolNames, expected.toolCount, "tool names should match count")

			for _, expectedTool := range expected.hasTools {
				assert.Contains(t, meta.ToolNames, expectedTool, "should include %s", expectedTool)
			}
		})
	}
}

func TestMetaToolNames(t *testing.T) {
	assert.Len(t, MetaToolNames, 4)
	assert.Contains(t, MetaToolNames, "help")
	assert.Contains(t, MetaToolNames, "get_tool_config")
	assert.Contains(t, MetaToolNames, "set_tool_config")
	assert.Contains(t, MetaToolNames, "list_tool_groups")
}

// TestToolCountConsistency ensures ToolCount field matches len(ToolNames) for all groups.
// This catches accidental desync when adding/removing tools from a group.
func TestToolCountConsistency(t *testing.T) {
	for _, groupName := range ToolGroupOrder {
		t.Run(groupName, func(t *testing.T) {
			meta := toolGroupMetadata[groupName]
			require.NotNil(t, meta, "metadata should exist for %s", groupName)

			assert.Equal(t, meta.ToolCount, len(meta.ToolNames),
				"ToolCount (%d) must match len(ToolNames) (%d) for group %s",
				meta.ToolCount, len(meta.ToolNames), groupName)
		})
	}
}

// TestToolGroupOrderConsistency ensures ToolGroupOrder matches toolGroupMetadata keys.
func TestToolGroupOrderConsistency(t *testing.T) {
	// Every group in ToolGroupOrder should exist in metadata
	for _, groupName := range ToolGroupOrder {
		_, exists := toolGroupMetadata[groupName]
		assert.True(t, exists, "group %s in ToolGroupOrder missing from toolGroupMetadata", groupName)
	}

	// Every group in metadata should be in ToolGroupOrder
	for groupName := range toolGroupMetadata {
		found := false
		for _, orderedName := range ToolGroupOrder {
			if orderedName == groupName {
				found = true
				break
			}
		}
		assert.True(t, found, "group %s in toolGroupMetadata missing from ToolGroupOrder", groupName)
	}

	// Counts should match
	assert.Equal(t, len(ToolGroupOrder), len(toolGroupMetadata),
		"ToolGroupOrder and toolGroupMetadata should have same number of entries")
}

// TestTotalToolCountMatchesDocumentation validates that tool counts in metadata
// sum to the documented total (72 tools = 68 group tools + 4 meta-tools).
// This catches drift between code and documentation.
func TestTotalToolCountMatchesDocumentation(t *testing.T) {
	// Sum all tool counts from metadata
	var groupToolCount int
	for _, meta := range toolGroupMetadata {
		groupToolCount += meta.ToolCount
	}

	// Add meta-tools
	totalTools := groupToolCount + len(MetaToolNames)

	// Expected total from documentation (CLAUDE.md, README.md, verify-docs.sh)
	const expectedTotal = 72

	assert.Equal(t, expectedTotal, totalTools,
		"Total tool count (%d group tools + %d meta-tools = %d) should match documented %d",
		groupToolCount, len(MetaToolNames), totalTools, expectedTotal)
}

// TestToolNamesAreUnique ensures no duplicate tool names exist across groups.
func TestToolNamesAreUnique(t *testing.T) {
	seen := make(map[string]string) // tool name -> group name

	for groupName, meta := range toolGroupMetadata {
		for _, toolName := range meta.ToolNames {
			if existingGroup, exists := seen[toolName]; exists {
				t.Errorf("Tool '%s' appears in both '%s' and '%s' groups",
					toolName, existingGroup, groupName)
			}
			seen[toolName] = groupName
		}
	}

	// Also check meta-tools don't conflict with group tools
	for _, metaTool := range MetaToolNames {
		if existingGroup, exists := seen[metaTool]; exists {
			t.Errorf("Meta-tool '%s' conflicts with tool in group '%s'",
				metaTool, existingGroup)
		}
	}
}
