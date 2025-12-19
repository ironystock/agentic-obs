// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

// Response type constants for UI messages.
const (
	// ResponseTypeReceived acknowledges receipt of an action.
	ResponseTypeReceived = "ui-message-received"
	// ResponseTypeResponse sends the result of processing an action.
	ResponseTypeResponse = "ui-message-response"
)

// UIResponse is sent from host to iframe in response to actions.
type UIResponse struct {
	// Type is the response type (ui-message-received or ui-message-response).
	Type string `json:"type"`
	// MessageID correlates this response to the originating action.
	MessageID string `json:"messageId"`
	// Payload contains the response data (only for ui-message-response).
	Payload *ResponsePayload `json:"payload,omitempty"`
}

// ResponsePayload is the payload structure for ui-message-response.
type ResponsePayload struct {
	// Response contains the successful result data.
	Response any `json:"response,omitempty"`
	// Error contains error information if the action failed.
	Error *ResponseError `json:"error,omitempty"`
}

// ResponseError contains error information for failed actions.
type ResponseError struct {
	// Message is a human-readable error description.
	Message string `json:"message"`
	// Code is an optional error code.
	Code string `json:"code,omitempty"`
	// Data contains additional error context.
	Data any `json:"data,omitempty"`
}

// NewReceivedResponse creates an acknowledgment response.
// This should be sent immediately when an action is received to confirm receipt.
func NewReceivedResponse(messageID string) *UIResponse {
	return &UIResponse{
		Type:      ResponseTypeReceived,
		MessageID: messageID,
	}
}

// NewSuccessResponse creates a success response with the result data.
// Use this when an action has been processed successfully.
func NewSuccessResponse(messageID string, result any) *UIResponse {
	return &UIResponse{
		Type:      ResponseTypeResponse,
		MessageID: messageID,
		Payload: &ResponsePayload{
			Response: result,
		},
	}
}

// NewErrorResponse creates an error response.
// Use this when an action fails to process.
func NewErrorResponse(messageID string, err error) *UIResponse {
	return &UIResponse{
		Type:      ResponseTypeResponse,
		MessageID: messageID,
		Payload: &ResponsePayload{
			Error: &ResponseError{
				Message: err.Error(),
			},
		},
	}
}

// NewErrorResponseWithCode creates an error response with an error code.
func NewErrorResponseWithCode(messageID string, code string, message string) *UIResponse {
	return &UIResponse{
		Type:      ResponseTypeResponse,
		MessageID: messageID,
		Payload: &ResponsePayload{
			Error: &ResponseError{
				Code:    code,
				Message: message,
			},
		},
	}
}

// NewErrorResponseWithData creates an error response with additional context data.
func NewErrorResponseWithData(messageID string, err error, data any) *UIResponse {
	return &UIResponse{
		Type:      ResponseTypeResponse,
		MessageID: messageID,
		Payload: &ResponsePayload{
			Error: &ResponseError{
				Message: err.Error(),
				Data:    data,
			},
		},
	}
}

// IsSuccess returns true if this response indicates success.
func (r *UIResponse) IsSuccess() bool {
	if r.Type == ResponseTypeReceived {
		return true
	}
	return r.Payload != nil && r.Payload.Error == nil
}

// IsError returns true if this response indicates an error.
func (r *UIResponse) IsError() bool {
	return r.Type == ResponseTypeResponse && r.Payload != nil && r.Payload.Error != nil
}

// GetError returns the error if present, nil otherwise.
func (r *UIResponse) GetError() *ResponseError {
	if r.Payload == nil {
		return nil
	}
	return r.Payload.Error
}

// GetResponse returns the response data if present, nil otherwise.
func (r *UIResponse) GetResponse() any {
	if r.Payload == nil {
		return nil
	}
	return r.Payload.Response
}
