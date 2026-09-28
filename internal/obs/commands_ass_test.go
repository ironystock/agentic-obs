package obs

import (
	"reflect"
	"testing"
)

// Most of the ASS surface lives inside a single function — the payload
// conversion — and three thin validation wrappers around CallVendorRequest.
// Because CallVendorRequest needs a live goobs connection, we exercise the
// validation paths against an unconnected Client (which returns an error on
// the first goobs call), and we test the pure payload helper directly.

func TestAssVariablesToVendorPayload(t *testing.T) {
	tests := []struct {
		name string
		in   []ASSVariable
		want []map[string]any
	}{
		{
			name: "empty",
			in:   []ASSVariable{},
			want: []map[string]any{},
		},
		{
			name: "single",
			in:   []ASSVariable{{Name: "current_game", Value: "Elden Ring"}},
			want: []map[string]any{{"name": "current_game", "value": "Elden Ring"}},
		},
		{
			name: "multiple preserves order",
			in: []ASSVariable{
				{Name: "a", Value: "1"},
				{Name: "b", Value: "2"},
				{Name: "c", Value: "3"},
			},
			want: []map[string]any{
				{"name": "a", "value": "1"},
				{"name": "b", "value": "2"},
				{"name": "c", "value": "3"},
			},
		},
		{
			name: "empty-string values are preserved",
			in:   []ASSVariable{{Name: "cleared", Value: ""}},
			want: []map[string]any{{"name": "cleared", "value": ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assVariablesToVendorPayload(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("payload mismatch.\n got: %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestASSSendMessageRejectsEmpty(t *testing.T) {
	c := NewClient(ConnectionConfig{Host: "127.0.0.1", Port: "4455"})
	err := c.ASSSendMessage("")
	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestASSRunMacroRejectsEmptyName(t *testing.T) {
	c := NewClient(ConnectionConfig{Host: "127.0.0.1", Port: "4455"})
	err := c.ASSRunMacro("", nil)
	if err == nil {
		t.Fatal("expected error for empty macro name, got nil")
	}
}

func TestASSSetVariablesRejectsEmpty(t *testing.T) {
	c := NewClient(ConnectionConfig{Host: "127.0.0.1", Port: "4455"})
	err := c.ASSSetVariables(nil)
	if err == nil {
		t.Fatal("expected error for empty variables list, got nil")
	}
}

// Non-validation paths require a live OBS+ASS to round-trip. Build-tag
// integration tests can cover those; the unit suite stops at validation.
