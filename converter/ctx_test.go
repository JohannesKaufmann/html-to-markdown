package converter

import (
	"context"
	"testing"
)

func TestState(t *testing.T) {
	state := newGlobalState()

	ctx := context.Background()
	ctx = state.provideGlobalState(ctx)

	val := GetState[int](ctx, "key")
	if val != 0 {
		t.Errorf("expected different value but got %d", val)
	}

	SetState[int](ctx, "key", 10)

	UpdateState[int](ctx, "key", func(i int) int {
		return i + 5
	})

	val = GetState[int](ctx, "key")
	if val != 15 {
		t.Errorf("expected different value but got %d", val)
	}
}

func TestContext(t *testing.T) {
	conv := NewConverter()
	bgCtx := context.Background()

	ctx := newConverterContext(bgCtx, conv)

	ctx1 := ctx.WithValue("keyA", "a1")
	if ctx1.Value("keyA") != "a1" {
		t.Error("got different value")
	}

	ctx2 := ctx.WithValue("keyA", "a2")
	if ctx2.Value("keyA") != "a2" {
		t.Error("got different value")
	}

	ctx3 := ctx.WithValue("keyB", "b1")
	if ctx3.Value("keyA") != nil {
		t.Error("expected nil value")
	}
	if ctx3.Value("keyB") != "b1" {
		t.Error("got different value")
	}
}
