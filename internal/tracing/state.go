package tracing

import (
	"context"
	"fmt"
	"sync"

	"claude-instrumentation/internal/hooks"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// hook events are not one single request, they arrive in HTTP calls.
// To build a trace tree from those disconnected events we need a state

type TraceState struct {
	mu sync.Mutex

	sessionsContexts map[string]context.Context
	sessionSpans     map[string]trace.Span
	turnContexts     map[string]context.Context
	turnSpans        map[string]trace.Span
	toolSpans        map[string]trace.Span
	agentContexts    map[string]context.Context
	agentSpans       map[string]trace.Span
}

func NewTraceState() *TraceState {
	return &TraceState{
		sessionsContexts: make(map[string]context.Context),
		sessionSpans:     make(map[string]trace.Span),
		turnContexts:     make(map[string]context.Context),
		turnSpans:        make(map[string]trace.Span),
		toolSpans:        make(map[string]trace.Span),
		agentContexts:    make(map[string]context.Context),
		agentSpans:       make(map[string]trace.Span),
	}
}

func (t *TraceState) SetSessionAttrs(sessionId string, context context.Context, span trace.Span) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.sessionsContexts[sessionId] = context
	t.sessionSpans[sessionId] = span
}

func (t *TraceState) GetSessionContext(sessionId string) (context.Context, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	ctx, ok := t.sessionsContexts[sessionId]
	return ctx, ok
}

func (t *TraceState) PopSessionSpan(sessionId string) (trace.Span, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, ok := t.sessionSpans[sessionId]
	if ok {
		delete(t.sessionSpans, sessionId)
		delete(t.sessionsContexts, sessionId)
	}
	return span, ok
}

func (t *TraceState) SetTurnAttrs(sessionId string, context context.Context, span trace.Span) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.turnContexts[sessionId] = context
	t.turnSpans[sessionId] = span
}

func (t *TraceState) GetTurnContext(sessionId string) (context.Context, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	ctx, ok := t.turnContexts[sessionId]
	return ctx, ok
}

func (t *TraceState) GetTurnSpan(sessionId string) (trace.Span, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, ok := t.turnSpans[sessionId]
	return span, ok
}

func (t *TraceState) PopTurnSpan(sessionId string) (trace.Span, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, ok := t.turnSpans[sessionId]
	if ok {
		delete(t.turnSpans, sessionId)
		delete(t.turnContexts, sessionId)
	}
	return span, ok
}

// Tool methods
func (t *TraceState) SetToolSpan(toolUseId string, span trace.Span) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.toolSpans[toolUseId] = span
}

func (t *TraceState) PopToolSpan(toolUseId string) (trace.Span, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, ok := t.toolSpans[toolUseId]
	if ok {
		delete(t.toolSpans, toolUseId)
	}
	return span, ok
}

// Agent methods
func (t *TraceState) SetAgentAttrs(agentId string, span trace.Span, context context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.agentSpans[agentId] = span
	t.agentContexts[agentId] = context
}

func (t *TraceState) GetAgentContext(agentId string) (context.Context, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	ctx, ok := t.agentContexts[agentId]
	return ctx, ok
}

func (t *TraceState) PopAgentSpan(agentId string) (trace.Span, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	span, ok := t.agentSpans[agentId]
	if ok {
		delete(t.agentSpans, agentId)
		delete(t.agentContexts, agentId)
	}
	return span, ok
}

func RecordHook(state *TraceState, input hooks.HookInput) error {
	sessionID := input.GetBaseInput().SessionID
	tracer := otel.Tracer("claude-instrumentation")

	switch input.GetEventName() {
	case hooks.EventSessionStart:
		sessionStart, ok := input.(hooks.SessionStartHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}

		ctx := context.Background()
		ctx, span := tracer.Start(ctx, fmt.Sprintf("session-%s", sessionID))
		span.SetAttributes(
			StringAttr("hook_event_name", string(sessionStart.HookEventName)),
			StringAttr("session_id", sessionStart.SessionID),
			StringAttr("cwd", sessionStart.CWD),
			StringAttr("model", sessionStart.Model),
			StringAttr("source", sessionStart.Source),
		)
		state.SetSessionAttrs(sessionID, ctx, span)

	case hooks.EventUserPromptSubmit:
		parentCtx, ok := state.GetSessionContext(sessionID)
		if !ok {
			parentCtx = context.Background()
		}

		userPromptSubmit, ok := input.(hooks.UserPromptSubmitHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}

		ctx, span := tracer.Start(parentCtx, fmt.Sprintf("turn-%s", sessionID))
		span.SetAttributes(
			StringAttr("hook_event_name", string(userPromptSubmit.HookEventName)),
			StringAttr("session_id", userPromptSubmit.SessionID),
			StringAttr("prompt", userPromptSubmit.Prompt),
		)
		state.SetTurnAttrs(sessionID, ctx, span)

	case hooks.EventPreToolUse:
		preToolUse, ok := input.(hooks.PreToolUseHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}

		parentCtx, ok := state.GetTurnContext(sessionID)
		if !ok {
			parentCtx = context.Background()
		}

		_, span := tracer.Start(parentCtx, fmt.Sprintf("tool-call-%s", preToolUse.ToolUseID))
		span.SetAttributes(
			StringAttr("hook_event_name", string(preToolUse.HookEventName)),
			StringAttr("session_id", preToolUse.SessionID),
			StringAttr("tool_name", preToolUse.ToolName),
			StringAttr("tool_use_id", preToolUse.ToolUseID),
		)
		state.SetToolSpan(preToolUse.ToolUseID, span)

	case hooks.EventPostToolUse:
		postToolUse, ok := input.(hooks.PostToolUseHookInput)
		if !ok {
			return hooks.ErrInvalidHookInput
		}
		span, ok := state.PopToolSpan(postToolUse.ToolUseID)
		if ok {
			span.SetAttributes(StringAttr("hook_event_name", string(postToolUse.HookEventName)))
			span.SetStatus(codes.Ok, "tool completed")
			span.End()
		}

	case hooks.EventStop:
		span, ok := state.PopTurnSpan(sessionID)
		if ok {
			span.SetStatus(codes.Ok, "turn completed")
			span.End()
		}

	case hooks.EventSessionEnd:
		span, ok := state.PopSessionSpan(sessionID)
		if ok {
			span.SetStatus(codes.Ok, "session completed")
			span.End()
		}
	}

	return nil
}
