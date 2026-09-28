package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Advanced Scene Switcher (ASS) tool input types.
//
// ASS exposes only three vendor request types over obs-websocket and they are
// all fire-and-forget — there is no introspection (no way to list macros or
// variables). Callers must know the names they want to use. The variable
// payload is documented as a list of {name, value} pairs where value is
// always a string on the wire.

// ASSVariableInput mirrors the obs.ASSVariable typed struct but accepts
// `any` for value so the LLM can pass numbers / bools naturally; the
// handler coerces to string before forwarding to OBS.
type ASSVariableInput struct {
	Name  string `json:"name" jsonschema:"Variable name (case-sensitive, must match ASS config)"`
	Value any    `json:"value" jsonschema:"Variable value (coerced to string)"`
}

// ASSRunMacroInput drives the ass_run_macro tool.
type ASSRunMacroInput struct {
	Name      string             `json:"name" jsonschema:"Exact name of the macro to run (case-sensitive)"`
	Variables []ASSVariableInput `json:"variables,omitempty" jsonschema:"Optional variables to set atomically before the macro runs"`
}

// ASSSendMessageInput drives the ass_send_message tool.
type ASSSendMessageInput struct {
	Message string `json:"message" jsonschema:"Message string to broadcast; macros with a matching websocket-message condition will react"`
}

// ASSSetVariablesInput drives the ass_set_variables tool (bulk).
type ASSSetVariablesInput struct {
	Variables []ASSVariableInput `json:"variables" jsonschema:"List of {name, value} pairs to set"`
}

// ASSSetVariableInput drives the ass_set_variable tool (single, ergonomic).
type ASSSetVariableInput struct {
	Name  string `json:"name" jsonschema:"Variable name (case-sensitive)"`
	Value any    `json:"value" jsonschema:"Variable value (coerced to string)"`
}

// coerceASSVariables converts the LLM-facing []ASSVariableInput (with any
// values) into the internal obs.ASSVariable shape (strings only). Numbers,
// bools, etc. become their fmt.Sprintf("%v", ...) form so users can write
// `{"name": "score", "value": 42}` without thinking about JSON typing.
func coerceASSVariables(in []ASSVariableInput) ([]obs.ASSVariable, error) {
	out := make([]obs.ASSVariable, 0, len(in))
	for i, v := range in {
		if v.Name == "" {
			return nil, fmt.Errorf("variables[%d].name must not be empty", i)
		}
		out = append(out, obs.ASSVariable{
			Name:  v.Name,
			Value: fmt.Sprintf("%v", v.Value),
		})
	}
	return out, nil
}

// handleASSRunMacro executes a named ASS macro, optionally pre-setting variables
// in the same atomic vendor call. This is the main natural-language entry point —
// most "trigger X in OBS" voice commands route here.
func (s *Server) handleASSRunMacro(ctx context.Context, request *mcpsdk.CallToolRequest, input ASSRunMacroInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("ASS: run macro %q (variables=%d)", input.Name, len(input.Variables))

	if input.Name == "" {
		return nil, nil, fmt.Errorf("name must not be empty")
	}

	vars, err := coerceASSVariables(input.Variables)
	if err != nil {
		return nil, nil, err
	}

	if err := s.obsClient.ASSRunMacro(input.Name, vars); err != nil {
		s.recordAction("ass_run_macro", "Run ASS macro", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to run ASS macro %q: %w", input.Name, err)
	}

	result := map[string]interface{}{
		"name":            input.Name,
		"variables_count": len(vars),
		"message":         fmt.Sprintf("Triggered ASS macro %q", input.Name),
	}
	s.recordAction("ass_run_macro", "Run ASS macro", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleASSSendMessage fires the Advanced Scene Switcher websocket-message
// event. Use this to trigger any macro configured with the "Websocket message
// received" condition.
func (s *Server) handleASSSendMessage(ctx context.Context, request *mcpsdk.CallToolRequest, input ASSSendMessageInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("ASS: send message %q", input.Message)

	if input.Message == "" {
		return nil, nil, fmt.Errorf("message must not be empty")
	}

	if err := s.obsClient.ASSSendMessage(input.Message); err != nil {
		s.recordAction("ass_send_message", "Send ASS message", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to send ASS message: %w", err)
	}

	result := map[string]interface{}{
		"message_sent": input.Message,
		"message":      fmt.Sprintf("Sent ASS websocket message %q", input.Message),
	}
	s.recordAction("ass_send_message", "Send ASS message", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleASSSetVariables bulk-updates ASS variables without firing any macro.
// Use when you want to update plugin state and let ASS's existing conditions
// react on their next evaluation pass.
func (s *Server) handleASSSetVariables(ctx context.Context, request *mcpsdk.CallToolRequest, input ASSSetVariablesInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("ASS: set %d variables", len(input.Variables))

	if len(input.Variables) == 0 {
		return nil, nil, fmt.Errorf("variables list must not be empty")
	}

	vars, err := coerceASSVariables(input.Variables)
	if err != nil {
		return nil, nil, err
	}

	if err := s.obsClient.ASSSetVariables(vars); err != nil {
		s.recordAction("ass_set_variables", "Set ASS variables", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set ASS variables: %w", err)
	}

	names := make([]string, len(vars))
	for i, v := range vars {
		names[i] = v.Name
	}
	result := map[string]interface{}{
		"count":   len(vars),
		"names":   names,
		"message": fmt.Sprintf("Set %d ASS variable(s)", len(vars)),
	}
	s.recordAction("ass_set_variables", "Set ASS variables", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleASSSetVariable is an ergonomic single-variable shortcut for the very
// common case where the LLM wants to update one variable by name. It routes
// to the same vendor request as the bulk handler.
func (s *Server) handleASSSetVariable(ctx context.Context, request *mcpsdk.CallToolRequest, input ASSSetVariableInput) (*mcpsdk.CallToolResult, any, error) {
	return s.handleASSSetVariables(ctx, request, ASSSetVariablesInput{
		Variables: []ASSVariableInput{{Name: input.Name, Value: input.Value}},
	})
}
