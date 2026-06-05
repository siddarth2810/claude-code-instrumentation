package tracing

import "testing"

func TestTraceStateStoresAndPopsSession(t *testing.T) {
	state := NewTraceState()

	state.SetSessionAttrs("session-1", "session-context-1", "session-span-1")

	ctx, ok := state.GetSessionContext("session-1")
	if !ok {
		t.Fatal("expected session context to be stored")
	}

	if ctx != "session-context-1" {
		t.Fatalf("session context = %q, want %q", ctx, "session-context-1")
	}

	span, ok := state.PopSessionSpan("session-1")
	if !ok {
		t.Fatal("expected session span to be popped")
	}

	if span != "session-span-1" {
		t.Fatalf("session span = %q, want %q", span, "session-span-1")
	}

	if _, ok := state.GetSessionContext("session-1"); ok {
		t.Fatal("session context should be removed after PopSessionSpan")
	}
}

func TestTraceStateStoresAndPopsTurnBySessionID(t *testing.T) {
	state := NewTraceState()

	state.SetTurnAttrs("session-1", "turn-context-1", "turn-span-1")

	ctx, ok := state.GetTurnContext("session-1")
	if !ok {
		t.Fatal("expected turn context to be stored")
	}

	if ctx != "turn-context-1" {
		t.Fatalf("turn context = %q, want %q", ctx, "turn-context-1")
	}

	span, ok := state.GetTurnSpan("session-1")
	if !ok {
		t.Fatal("expected turn span to be stored")
	}

	if span != "turn-span-1" {
		t.Fatalf("turn span = %q, want %q", span, "turn-span-1")
	}

	poppedSpan, ok := state.PopTurnSpan("session-1")
	if !ok {
		t.Fatal("expected turn span to be popped")
	}

	if poppedSpan != "turn-span-1" {
		t.Fatalf("popped turn span = %q, want %q", poppedSpan, "turn-span-1")
	}

	if _, ok := state.GetTurnContext("session-1"); ok {
		t.Fatal("turn context should be removed after PopTurnSpan")
	}
}

func TestTraceStateStoresAndPopsToolByToolUseID(t *testing.T) {
	state := NewTraceState()

	state.SetToolSpan("toolu-1", "tool-span-1")

	span, ok := state.PopToolSpan("toolu-1")
	if !ok {
		t.Fatal("expected tool span to be popped")
	}

	if span != "tool-span-1" {
		t.Fatalf("tool span = %q, want %q", span, "tool-span-1")
	}

	if _, ok := state.PopToolSpan("toolu-1"); ok {
		t.Fatal("tool span should not exist after first pop")
	}
}

func TestTraceStateStoresAgentSpanAndContextByAgentID(t *testing.T) {
	state := NewTraceState()

	state.SetAgentAttrs("agent-1", "agent-span-1", "agent-context-1")

	ctx, ok := state.GetAgentContext("agent-1")
	if !ok {
		t.Fatal("expected agent context to be stored")
	}

	if ctx != "agent-context-1" {
		t.Fatalf("agent context = %q, want %q", ctx, "agent-context-1")
	}

	span, ok := state.PopAgentSpan("agent-1")
	if !ok {
		t.Fatal("expected agent span to be popped")
	}

	if span != "agent-span-1" {
		t.Fatalf("agent span = %q, want %q", span, "agent-span-1")
	}

	if _, ok := state.GetAgentContext("agent-1"); ok {
		t.Fatal("agent context should be removed after PopAgentSpan")
	}
}
