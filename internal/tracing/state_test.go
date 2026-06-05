package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func testSpan() trace.Span {
	return trace.SpanFromContext(context.Background())
}

func TestTraceStateStoresAndPopsSession(t *testing.T) {
	state := NewTraceState()
	ctx := context.Background()
	span := testSpan()

	state.SetSessionAttrs("session-1", ctx, span)

	gotCtx, ok := state.GetSessionContext("session-1")
	if !ok {
		t.Fatal("expected session context to be stored")
	}

	if gotCtx != ctx {
		t.Fatal("session context did not match stored context")
	}

	gotSpan, ok := state.PopSessionSpan("session-1")
	if !ok {
		t.Fatal("expected session span to be popped")
	}

	if gotSpan != span {
		t.Fatal("session span did not match stored span")
	}

	if _, ok := state.GetSessionContext("session-1"); ok {
		t.Fatal("session context should be removed after PopSessionSpan")
	}
}

func TestTraceStateStoresAndPopsTurnBySessionID(t *testing.T) {
	state := NewTraceState()
	ctx := context.Background()
	span := testSpan()

	state.SetTurnAttrs("session-1", ctx, span)

	gotCtx, ok := state.GetTurnContext("session-1")
	if !ok {
		t.Fatal("expected turn context to be stored")
	}

	if gotCtx != ctx {
		t.Fatal("turn context did not match stored context")
	}

	gotSpan, ok := state.GetTurnSpan("session-1")
	if !ok {
		t.Fatal("expected turn span to be stored")
	}

	if gotSpan != span {
		t.Fatal("turn span did not match stored span")
	}

	poppedSpan, ok := state.PopTurnSpan("session-1")
	if !ok {
		t.Fatal("expected turn span to be popped")
	}

	if poppedSpan != span {
		t.Fatal("popped turn span did not match stored span")
	}

	if _, ok := state.GetTurnContext("session-1"); ok {
		t.Fatal("turn context should be removed after PopTurnSpan")
	}
}

func TestTraceStateStoresAndPopsToolByToolUseID(t *testing.T) {
	state := NewTraceState()
	span := testSpan()

	state.SetToolSpan("toolu-1", span)

	gotSpan, ok := state.PopToolSpan("toolu-1")
	if !ok {
		t.Fatal("expected tool span to be popped")
	}

	if gotSpan != span {
		t.Fatal("tool span did not match stored span")
	}

	if _, ok := state.PopToolSpan("toolu-1"); ok {
		t.Fatal("tool span should not exist after first pop")
	}
}

func TestTraceStateStoresAgentSpanAndContextByAgentID(t *testing.T) {
	state := NewTraceState()
	ctx := context.Background()
	span := testSpan()

	state.SetAgentAttrs("agent-1", span, ctx)

	gotCtx, ok := state.GetAgentContext("agent-1")
	if !ok {
		t.Fatal("expected agent context to be stored")
	}

	if gotCtx != ctx {
		t.Fatal("agent context did not match stored context")
	}

	gotSpan, ok := state.PopAgentSpan("agent-1")
	if !ok {
		t.Fatal("expected agent span to be popped")
	}

	if gotSpan != span {
		t.Fatal("agent span did not match stored span")
	}

	if _, ok := state.GetAgentContext("agent-1"); ok {
		t.Fatal("agent context should be removed after PopAgentSpan")
	}
}
