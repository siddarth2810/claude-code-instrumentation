package hooks

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
