// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReceivedResponse(t *testing.T) {
	resp := NewReceivedResponse("msg-123")

	assert.Equal(t, ResponseTypeReceived, resp.Type)
	assert.Equal(t, "msg-123", resp.MessageID)
	assert.Nil(t, resp.Payload)
	assert.True(t, resp.IsSuccess())
	assert.False(t, resp.IsError())
}

func TestNewSuccessResponse(t *testing.T) {
	result := map[string]any{
		"status": "ok",
		"data":   []string{"a", "b", "c"},
	}
	resp := NewSuccessResponse("msg-456", result)

	assert.Equal(t, ResponseTypeResponse, resp.Type)
	assert.Equal(t, "msg-456", resp.MessageID)
	assert.NotNil(t, resp.Payload)
	assert.Nil(t, resp.Payload.Error)
	assert.Equal(t, result, resp.Payload.Response)
	assert.True(t, resp.IsSuccess())
	assert.False(t, resp.IsError())
	assert.Equal(t, result, resp.GetResponse())
	assert.Nil(t, resp.GetError())
}

func TestNewErrorResponse(t *testing.T) {
	err := errors.New("something went wrong")
	resp := NewErrorResponse("msg-789", err)

	assert.Equal(t, ResponseTypeResponse, resp.Type)
	assert.Equal(t, "msg-789", resp.MessageID)
	assert.NotNil(t, resp.Payload)
	assert.NotNil(t, resp.Payload.Error)
	assert.Equal(t, "something went wrong", resp.Payload.Error.Message)
	assert.Empty(t, resp.Payload.Error.Code)
	assert.False(t, resp.IsSuccess())
	assert.True(t, resp.IsError())
	assert.Nil(t, resp.GetResponse())
	assert.NotNil(t, resp.GetError())
}

func TestNewErrorResponseWithCode(t *testing.T) {
	resp := NewErrorResponseWithCode("msg-abc", "INVALID_PARAMS", "Missing required parameter")

	assert.Equal(t, ResponseTypeResponse, resp.Type)
	assert.Equal(t, "msg-abc", resp.MessageID)
	assert.NotNil(t, resp.Payload.Error)
	assert.Equal(t, "INVALID_PARAMS", resp.Payload.Error.Code)
	assert.Equal(t, "Missing required parameter", resp.Payload.Error.Message)
}

func TestNewErrorResponseWithData(t *testing.T) {
	err := errors.New("validation failed")
	data := map[string]any{
		"field":   "email",
		"details": "invalid format",
	}
	resp := NewErrorResponseWithData("msg-def", err, data)

	assert.Equal(t, ResponseTypeResponse, resp.Type)
	assert.Equal(t, "msg-def", resp.MessageID)
	assert.NotNil(t, resp.Payload.Error)
	assert.Equal(t, "validation failed", resp.Payload.Error.Message)
	assert.Equal(t, data, resp.Payload.Error.Data)
}

func TestUIResponse_JSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		resp *UIResponse
	}{
		{
			name: "received response",
			resp: NewReceivedResponse("test-id"),
		},
		{
			name: "success response",
			resp: NewSuccessResponse("test-id", map[string]string{"key": "value"}),
		},
		{
			name: "error response",
			resp: NewErrorResponse("test-id", errors.New("test error")),
		},
		{
			name: "error with code",
			resp: NewErrorResponseWithCode("test-id", "ERR_001", "test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.resp)
			require.NoError(t, err)

			var decoded UIResponse
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			assert.Equal(t, tt.resp.Type, decoded.Type)
			assert.Equal(t, tt.resp.MessageID, decoded.MessageID)
		})
	}
}

func TestUIResponse_JSONFormat(t *testing.T) {
	t.Run("received format", func(t *testing.T) {
		resp := NewReceivedResponse("msg-123")
		data, err := json.Marshal(resp)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "ui-message-received", m["type"])
		assert.Equal(t, "msg-123", m["messageId"])
		_, hasPayload := m["payload"]
		assert.False(t, hasPayload, "payload should be omitted for received")
	})

	t.Run("success format", func(t *testing.T) {
		resp := NewSuccessResponse("msg-456", "result data")
		data, err := json.Marshal(resp)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "ui-message-response", m["type"])
		assert.Equal(t, "msg-456", m["messageId"])
		payload := m["payload"].(map[string]any)
		assert.Equal(t, "result data", payload["response"])
		_, hasError := payload["error"]
		assert.False(t, hasError)
	})

	t.Run("error format", func(t *testing.T) {
		resp := NewErrorResponse("msg-789", errors.New("test error"))
		data, err := json.Marshal(resp)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "ui-message-response", m["type"])
		payload := m["payload"].(map[string]any)
		errObj := payload["error"].(map[string]any)
		assert.Equal(t, "test error", errObj["message"])
	})
}

func TestResponseTypeConstants(t *testing.T) {
	// Verify constants match protocol specification
	assert.Equal(t, "ui-message-received", ResponseTypeReceived)
	assert.Equal(t, "ui-message-response", ResponseTypeResponse)
}

func TestUIResponse_HelperMethods(t *testing.T) {
	t.Run("received response helpers", func(t *testing.T) {
		resp := NewReceivedResponse("id")
		assert.True(t, resp.IsSuccess())
		assert.False(t, resp.IsError())
		assert.Nil(t, resp.GetError())
		assert.Nil(t, resp.GetResponse())
	})

	t.Run("success response helpers", func(t *testing.T) {
		resp := NewSuccessResponse("id", "data")
		assert.True(t, resp.IsSuccess())
		assert.False(t, resp.IsError())
		assert.Nil(t, resp.GetError())
		assert.Equal(t, "data", resp.GetResponse())
	})

	t.Run("error response helpers", func(t *testing.T) {
		resp := NewErrorResponse("id", errors.New("err"))
		assert.False(t, resp.IsSuccess())
		assert.True(t, resp.IsError())
		assert.NotNil(t, resp.GetError())
		assert.Nil(t, resp.GetResponse())
	})

	t.Run("nil payload response", func(t *testing.T) {
		resp := &UIResponse{Type: ResponseTypeResponse, MessageID: "id"}
		assert.False(t, resp.IsSuccess())
		assert.False(t, resp.IsError())
		assert.Nil(t, resp.GetError())
		assert.Nil(t, resp.GetResponse())
	})
}
