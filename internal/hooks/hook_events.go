package hooks

import (
	"encoding/json"
	"errors"
)

type EventName string

var ErrInvalidHookInput = errors.New("invalid hook input")

const (
	EventPreToolUse         EventName = "PreToolUse"
	EventSessionStart       EventName = "SessionStart"
	EventPostToolUse        EventName = "PostToolUse"
	EventPostToolUseFailure EventName = "PostToolUseFailure"
	EventUserPromptSubmit   EventName = "UserPromptSubmit"
	EventStop               EventName = "Stop"
	EventSubagentStop       EventName = "SubagentStop"
	EventPreCompact         EventName = "PreCompact"
	EventNotification       EventName = "Notification"
	EventSubagentStart      EventName = "SubagentStart"
	EventPermissionRequest  EventName = "PermissionRequest"
	EventStopFailure        EventName = "StopFailure"
	EventSessionEnd         EventName = "SessionEnd"
)

type SubagentContext struct {
	AgentID   string `json:"agent_id,omitempty"`
	AgentType string `json:"agent_type,omitempty"`
}

// Hook input types
type BaseHookInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	PermissionMode string `json:"permission_mode,omitempty"`
}

// Session start hook event
type SessionStartHookInput struct {
	BaseHookInput
	HookEventName EventName `json:"hook_event_name"`
	Source        string    `json:"source"`
	Model         string    `json:"model"`
}

// Input data for PreToolUse hook events
type PreToolUseHookInput struct {
	BaseHookInput
	HookEventName EventName `json:"hook_event_name"`
	ToolName      string    `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
	ToolUseID     string    `json:"tool_use_id"`
}

// PostToolUseHookInput is the input for PostToolUse hook events.
// ToolResponse shape also depends on ToolName so kept as raw JSON.
type PostToolUseHookInput struct {
	BaseHookInput
	SubagentContext
	HookEventName EventName       `json:"hook_event_name"`
	ToolName      string          `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
	ToolResponse  json.RawMessage `json:"tool_response"`
	ToolUseID     string          `json:"tool_use_id"`
}

// PostToolUseFailureHookInput is the input for PostToolUseFailure hook events.
type PostToolUseFailureHookInput struct {
	BaseHookInput
	SubagentContext
	HookEventName EventName       `json:"hook_event_name"`
	ToolName      string          `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
	ToolUseID     string          `json:"tool_use_id"`
	Error         string          `json:"error"`
	IsInterrupt   bool            `json:"is_interrupt,omitempty"` // NotRequired in Python
}

// UserPromptSubmitHookInput is the input for UserPromptSubmit hook events.
type UserPromptSubmitHookInput struct {
	BaseHookInput
	HookEventName EventName `json:"hook_event_name"`
	Prompt        string    `json:"prompt"`
}

// StopHookInput is the input for Stop hook events.
type StopHookInput struct {
	BaseHookInput
	HookEventName  EventName `json:"hook_event_name"`
	StopHookActive bool      `json:"stop_hook_active"`
}

// StopFailureHookInput is the input for StopFailure hook events.
type StopFailureHookInput struct {
	BaseHookInput
	HookEventName        EventName `json:"hook_event_name"`
	Error                string    `json:"error"`
	ErrorDetails         string    `json:"error_details,omitempty"`          // NotRequired
	LastAssistantMessage string    `json:"last_assistant_message,omitempty"` // NotRequired
}

// SessionEndHookInput is the input for SessionEnd hook events.
type SessionEndHookInput struct {
	BaseHookInput
	HookEventName EventName `json:"hook_event_name"`
	Reason        string    `json:"reason"`
}

// SubagentStopHookInput is the input for SubagentStop hook events.
// AgentID and AgentType are required here (unlike SubagentContext where they are optional).
type SubagentStopHookInput struct {
	BaseHookInput
	HookEventName       EventName `json:"hook_event_name"`
	StopHookActive      bool      `json:"stop_hook_active"`
	AgentID             string    `json:"agent_id"`
	AgentTranscriptPath string    `json:"agent_transcript_path"`
	AgentType           string    `json:"agent_type"`
}

// PreCompactHookInput is the input for PreCompact hook events.
type PreCompactHookInput struct {
	BaseHookInput
	HookEventName      EventName `json:"hook_event_name"`
	Trigger            string    `json:"trigger"`             // "manual" | "auto"
	CustomInstructions string    `json:"custom_instructions"` // empty string when null
}

// NotificationHookInput is the input for Notification hook events.
type NotificationHookInput struct {
	BaseHookInput
	HookEventName    EventName `json:"hook_event_name"`
	Message          string    `json:"message"`
	Title            string    `json:"title,omitempty"` // NotRequired
	NotificationType string    `json:"notification_type"`
}

// SubagentStartHookInput is the input for SubagentStart hook events.
// AgentID and AgentType are required (not optional like SubagentContext).
type SubagentStartHookInput struct {
	BaseHookInput
	HookEventName EventName `json:"hook_event_name"`
	AgentID       string    `json:"agent_id"`
	AgentType     string    `json:"agent_type"`
}

// PermissionRequestHookInput is the input for PermissionRequest hook events.
// PermissionSuggestions is kept as raw JSON — the shape is complex and
// you likely only need it for pass-through or logging.
type PermissionRequestHookInput struct {
	BaseHookInput
	SubagentContext
	HookEventName         EventName       `json:"hook_event_name"`
	ToolName              string          `json:"tool_name"`
	ToolInput             json.RawMessage `json:"tool_input"`
	PermissionSuggestions json.RawMessage `json:"permission_suggestions,omitempty"` // NotRequired
}


