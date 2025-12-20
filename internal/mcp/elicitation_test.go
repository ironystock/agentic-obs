package mcp

import (
	"context"
	"testing"
)

func TestElicitConfirmation_NilSession(t *testing.T) {
	// When session is nil, confirmation should be skipped (return true)
	ctx := context.Background()
	confirmed, err := ElicitConfirmation(ctx, nil, "Test message")

	if err != nil {
		t.Errorf("ElicitConfirmation(nil session) error = %v, want nil", err)
	}
	if !confirmed {
		t.Error("ElicitConfirmation(nil session) should return true when session is nil")
	}
}

func TestElicitStreamingConfirmation_NilSession(t *testing.T) {
	ctx := context.Background()
	confirmed, err := ElicitStreamingConfirmation(ctx, nil)

	if err != nil {
		t.Errorf("ElicitStreamingConfirmation(nil session) error = %v, want nil", err)
	}
	if !confirmed {
		t.Error("ElicitStreamingConfirmation(nil session) should return true when session is nil")
	}
}

func TestElicitStopStreamingConfirmation_NilSession(t *testing.T) {
	ctx := context.Background()
	confirmed, err := ElicitStopStreamingConfirmation(ctx, nil)

	if err != nil {
		t.Errorf("ElicitStopStreamingConfirmation(nil session) error = %v, want nil", err)
	}
	if !confirmed {
		t.Error("ElicitStopStreamingConfirmation(nil session) should return true when session is nil")
	}
}

func TestElicitDeleteConfirmation_NilSession(t *testing.T) {
	ctx := context.Background()
	confirmed, err := ElicitDeleteConfirmation(ctx, nil, "scene", "TestScene")

	if err != nil {
		t.Errorf("ElicitDeleteConfirmation(nil session) error = %v, want nil", err)
	}
	if !confirmed {
		t.Error("ElicitDeleteConfirmation(nil session) should return true when session is nil")
	}
}

func TestCancelledResult(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{"Streaming", "Streaming cancelled by user"},
		{"Scene deletion", "Scene deletion cancelled by user"},
		{"Recording", "Recording cancelled by user"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			result := CancelledResult(tt.action)
			if result.Message != tt.want {
				t.Errorf("CancelledResult(%q) = %q, want %q", tt.action, result.Message, tt.want)
			}
		})
	}
}
