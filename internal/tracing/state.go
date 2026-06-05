package tracing

import "claude-instrumentation/internal/hooks"

// hook events are not one single request, they arrive in HTTP calls.
// To build a trace tree from those disconnected events we need a state

type TraceState struct {
	sessionsContexts map[string]string
	sessionSpans     map[string]string
	turnContexts     map[string]string
	turnSpans        map[string]string
	toolSpans        map[string]string
	agentContexts    map[string]string
	agentSpans       map[string]string
}

func NewTraceState() *TraceState {
	return &TraceState{
		sessionsContexts: make(map[string]string),
		sessionSpans:     make(map[string]string),
		turnContexts:     make(map[string]string),
		turnSpans:        make(map[string]string),
		toolSpans:        make(map[string]string),
		agentContexts:    make(map[string]string),
		agentSpans:       make(map[string]string),
	}
}

func (t *TraceState) SetSessionAttrs(sessionId, context, span string) {
	t.sessionsContexts[sessionId] = context
	t.sessionSpans[sessionId] = span
}

func (t *TraceState) GetSessionContext(sessionId string) (string, bool) {
	ctx, ok := t.sessionsContexts[sessionId]
	return ctx, ok
}

func (t *TraceState) PopSessionSpan(sessionId string) (string, bool) {
	span, ok := t.sessionSpans[sessionId]
	if ok {
		delete(t.sessionSpans, sessionId)
		delete(t.sessionsContexts, sessionId)
	}
	return span, ok
}

func (t *TraceState) SetTurnAttrs(sessionId, context, span string) {
	t.turnContexts[sessionId] = context
	t.turnSpans[sessionId] = span
}

func (t *TraceState) GetTurnContext(sessionId string) (string, bool) {
	ctx, ok := t.turnContexts[sessionId]
	return ctx, ok
}

func (t *TraceState) GetTurnSpan(sessionId string) (string, bool) {
	span, ok := t.turnSpans[sessionId]
	return span, ok
}

func (t *TraceState) PopTurnSpan(sessionId string) (string, bool) {
	span, ok := t.turnSpans[sessionId]
	if ok {
		delete(t.turnSpans, sessionId)
		delete(t.turnContexts, sessionId)
	}
	return span, ok
}

// Tool methods
func (t *TraceState) SetToolSpan(toolUseId, span string) {
	t.toolSpans[toolUseId] = span
}

func (t *TraceState) PopToolSpan(toolUseId string) (string, bool) {
	span, ok := t.toolSpans[toolUseId]
	if ok {
		delete(t.toolSpans, toolUseId)
	}
	return span, ok
}

// Agent methods
func (t *TraceState) SetAgentAttrs(agentId, span, context string) {
	t.agentSpans[agentId] = span
	t.agentContexts[agentId] = context
}

func (t *TraceState) GetAgentContext(agentId string) (string, bool) {
	ctx, ok := t.agentContexts[agentId]
	return ctx, ok
}

func (t *TraceState) PopAgentSpan(agentId string) (string, bool) {
	span, ok := t.agentSpans[agentId]
	if ok {
		delete(t.agentSpans, agentId)
		delete(t.agentContexts, agentId)
	}
	return span, ok
}

func RecordHook(state *TraceState, input hooks.HookInput) error {
	sessionID := input.GetBaseInput().SessionID

	switch input.GetEventName() {
	case hooks.EventSessionStart:
		state.SetSessionAttrs(sessionID, "session-context", "session-span")

	case hooks.EventUserPromptSubmit:
		state.SetTurnAttrs(sessionID, "turn-context", "turn-span")

	case hooks.EventPreToolUse:
		preToolUse, ok := input.(hooks.PreToolUseHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}
		state.SetToolSpan(preToolUse.ToolUseID, "tool-span")

	case hooks.EventPostToolUse:
		postToolUse, ok := input.(hooks.PostToolUseHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}
		state.PopToolSpan(postToolUse.ToolUseID)

	case hooks.EventStop:
		state.PopTurnSpan(sessionID)

	case hooks.EventSessionEnd:
		state.PopSessionSpan(sessionID)
	}

	return nil
}
