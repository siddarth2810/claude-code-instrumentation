package hooks

import (
	"encoding/json"
)

type HookInput interface {
	GetEventName() EventName
	GetBaseInput() BaseHookInput
}

// Implement the interface for all structs by adding these helper methods.

func (h SessionStartHookInput) GetEventName() EventName     { return h.HookEventName }
func (h SessionStartHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h PreToolUseHookInput) GetEventName() EventName     { return h.HookEventName }
func (h PreToolUseHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h PostToolUseHookInput) GetEventName() EventName     { return h.HookEventName }
func (h PostToolUseHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h PostToolUseFailureHookInput) GetEventName() EventName     { return h.HookEventName }
func (h PostToolUseFailureHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h UserPromptSubmitHookInput) GetEventName() EventName     { return h.HookEventName }
func (h UserPromptSubmitHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h StopHookInput) GetEventName() EventName     { return h.HookEventName }
func (h StopHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h StopFailureHookInput) GetEventName() EventName     { return h.HookEventName }
func (h StopFailureHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h SessionEndHookInput) GetEventName() EventName     { return h.HookEventName }
func (h SessionEndHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h SubagentStopHookInput) GetEventName() EventName     { return h.HookEventName }
func (h SubagentStopHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h PreCompactHookInput) GetEventName() EventName     { return h.HookEventName }
func (h PreCompactHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h NotificationHookInput) GetEventName() EventName     { return h.HookEventName }
func (h NotificationHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h SubagentStartHookInput) GetEventName() EventName     { return h.HookEventName }
func (h SubagentStartHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func (h PermissionRequestHookInput) GetEventName() EventName     { return h.HookEventName }
func (h PermissionRequestHookInput) GetBaseInput() BaseHookInput { return h.BaseHookInput }

func DecodeHookInput(body []byte) (HookInput, error) {
	var event struct {
		HookEventName EventName `json:"hook_event_name"`
	}

	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, err
	}

	if event.HookEventName == "" {
		return nil, ErrInvalidHookInput
	}

	switch event.HookEventName {
	case EventSessionStart:
		var input SessionStartHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventPreToolUse:
		var input PreToolUseHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventPostToolUse:
		var input PostToolUseHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventPostToolUseFailure:
		var input PostToolUseFailureHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventUserPromptSubmit:
		var input UserPromptSubmitHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventStop:
		var input StopHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventSubagentStop:
		var input SubagentStopHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventPreCompact:
		var input PreCompactHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventNotification:
		var input NotificationHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventSubagentStart:
		var input SubagentStartHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventPermissionRequest:
		var input PermissionRequestHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventStopFailure:
		var input StopFailureHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	case EventSessionEnd:
		var input SessionEndHookInput
		if err := json.Unmarshal(body, &input); err != nil {
			return nil, err
		}
		return input, nil
	default:
		return nil, ErrInvalidHookInput
	}
}
