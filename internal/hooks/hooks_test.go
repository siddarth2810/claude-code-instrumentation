package hooks

import (
	"errors"
	"testing"
)

func TestDecodeHookInputDecodesSessionStart(t *testing.T) {
	body := []byte(`{
		"session_id": "session-123",
		"transcript_path": "/tmp/transcript.jsonl",
		"cwd": "/home/sid/project",
		"hook_event_name": "SessionStart",
		"source": "startup",
		"model": "claude-sonnet-4"
	}`)

	got, err := DecodeHookInput(body)
	if err != nil {
		t.Fatalf("DecodeHookInput returned error: %v", err)
	}

	event, ok := got.(SessionStartHookInput)
	if !ok {
		t.Fatalf("DecodeHookInput returned %T, want SessionStartHookInput", got)
	}

	if event.GetEventName() != EventSessionStart {
		t.Fatalf("event name = %q, want %q", event.GetEventName(), EventSessionStart)
	}

	if event.SessionID != "session-123" {
		t.Fatalf("session ID = %q, want %q", event.SessionID, "session-123")
	}

	if event.Source != "startup" {
		t.Fatalf("source = %q, want %q", event.Source, "startup")
	}

	if event.Model != "claude-sonnet-4" {
		t.Fatalf("model = %q, want %q", event.Model, "claude-sonnet-4")
	}
}

func TestDecodeHookInputDecodesPreToolUse(t *testing.T) {
	body := []byte(`{
		"session_id": "session-123",
		"transcript_path": "/tmp/transcript.jsonl",
		"cwd": "/home/sid/project",
		"hook_event_name": "PreToolUse",
		"tool_name": "Bash",
		"tool_input": {"command": "go test ./..."},
		"tool_use_id": "toolu-123"
	}`)

	got, err := DecodeHookInput(body)
	if err != nil {
		t.Fatalf("DecodeHookInput returned error: %v", err)
	}

	event, ok := got.(PreToolUseHookInput)
	if !ok {
		t.Fatalf("DecodeHookInput returned %T, want PreToolUseHookInput", got)
	}

	if event.GetEventName() != EventPreToolUse {
		t.Fatalf("event name = %q, want %q", event.GetEventName(), EventPreToolUse)
	}

	if event.ToolName != "Bash" {
		t.Fatalf("tool name = %q, want %q", event.ToolName, "Bash")
	}

	if event.ToolUseID != "toolu-123" {
		t.Fatalf("tool use ID = %q, want %q", event.ToolUseID, "toolu-123")
	}

	if len(event.ToolInput) == 0 {
		t.Fatal("tool input was empty, want raw JSON to be preserved")
	}
}

func TestDecodeHookInputRejectsMissingEventName(t *testing.T) {
	body := []byte(`{
		"session_id": "session-123",
		"transcript_path": "/tmp/transcript.jsonl",
		"cwd": "/home/sid/project"
	}`)

	got, err := DecodeHookInput(body)
	if !errors.Is(err, ErrInvalidHookInput) {
		t.Fatalf("DecodeHookInput error = %v, want ErrInvalidHookInput", err)
	}

	if got != nil {
		t.Fatalf("DecodeHookInput returned %T, want nil input", got)
	}
}

func TestDecodeHookInputRejectsUnknownEventName(t *testing.T) {
	body := []byte(`{
		"session_id": "session-123",
		"transcript_path": "/tmp/transcript.jsonl",
		"cwd": "/home/sid/project",
		"hook_event_name": "SomethingElse"
	}`)

	got, err := DecodeHookInput(body)
	if !errors.Is(err, ErrInvalidHookInput) {
		t.Fatalf("DecodeHookInput error = %v, want ErrInvalidHookInput", err)
	}

	if got != nil {
		t.Fatalf("DecodeHookInput returned %T, want nil input", got)
	}
}
