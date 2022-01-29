package callbacks

import (
	"testing"
)

func TestCallbacks_New_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	cbs = New("Test", func(Key) error { return nil }, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	cbs = New("Test", nil, func(Key) error { return nil })
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	cbs = New("Test", func(Key) error { return nil },
		func(Key) error { return nil })
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}
}

func TestCallbacks_AddSingleShot_Errors(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddSingleShot(Key(0), nil)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	_, err = cbs.AddSingleShot(Key(0), "non-function")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_AddSingleShot_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddSingleShot(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
}

func TestCallbacks_AddSingleShot_FirstFunc_Success(t *testing.T) {
	i := 0

	cbs := New("Test", func(Key) error { i++; return nil }, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddSingleShot(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	_, err = cbs.AddSingleShot(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	if i != 1 {
		t.Fatalf("expected firstFunc to be called once, got called %d "+
			"times", i)
	}
}

func TestCallbacks_AddContinuous_Errors(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddContinuous(Key(0), nil)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	_, err = cbs.AddContinuous(Key(0), "non-function")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_AddContinuous_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
}

func TestCallbacks_AddContinuous_FirstFunc_Success(t *testing.T) {
	i := 0

	cbs := New("Test", func(Key) error { i++; return nil }, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	if i != 1 {
		t.Fatalf("expected firstFunc to be called once, got called %d "+
			"times", i)
	}

	_, err = cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	if i != 1 {
		t.Fatalf("expected firstFunc to be called once, got called %d "+
			"times", i)
	}
}

func TestCallbacks_Remove_Errors(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	err := cbs.Remove(Key(0), Tag(0)) // Invalid tag.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	err = cbs.Remove(Key(0), Tag(1)) // Key not found.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	tag, err := cbs.AddSingleShot(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	err = cbs.Remove(Key(0), Tag(2)) // Tag not found for key.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	err = cbs.Remove(Key(0), tag) // Can not remove single-shot callback.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_Remove_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	tag, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	err = cbs.Remove(Key(0), tag)
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	err = cbs.Remove(Key(0), tag)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_Remove_LastFunc_Success(t *testing.T) {
	i := 0

	cbs := New("Test", nil, func(Key) error { i++; return nil })
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	tag1, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	tag2, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	err = cbs.Remove(Key(0), tag1)
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	if i != 0 {
		t.Fatalf("expected lastFunc not to be called, got called %d "+
			"times", i)
	}

	err = cbs.Remove(Key(0), tag2)
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	if i != 1 {
		t.Fatalf("expected lastFunc to be called once, got called %d "+
			"times", i)
	}
}

func TestCallbacks_Callback_Errors(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.Callback(Key(0), Tag(0)) // Invalid tag.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	_, err = cbs.Callback(Key(0), Tag(1)) // Key not found.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	tag, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	_, err = cbs.Callback(Key(0), Tag(tag+1)) // Tag not found for key.
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_Continuous_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	i := 0

	tag, err := cbs.AddContinuous(Key(0), func() { i++ })
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	cb, err := cbs.Callback(Key(0), Tag(tag))
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	if cb == nil {
		t.Fatalf("expected non-nil callback, got nil")
	}

	cbFunc, ok := cb.(func())
	if !ok {
		t.Fatalf("expected func(), got something %#+v", cb)
	}

	cbFunc()

	if i != 1 {
		t.Fatalf("expected callback to be called once, got called %d "+
			"times", i)

	}

	cb, err = cbs.Callback(Key(0), Tag(tag))
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	if cb == nil {
		t.Fatalf("expected non-nil callback, got nil")
	}
}

func TestCallbacks_SingleShot_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	i := 0

	tag, err := cbs.AddSingleShot(Key(0), func() { i++ })
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	cb, err := cbs.Callback(Key(0), tag)
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}
	if cb == nil {
		t.Fatalf("expected non-nil callback, got nil")
	}

	cbFunc, ok := cb.(func())
	if !ok {
		t.Fatalf("expected func(), got something %#+v", cb)
	}

	cbFunc()

	if i != 1 {
		t.Fatalf("expected callback to be called once, got called %d "+
			"times", i)

	}

	cb, err = cbs.Callback(Key(0), tag)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_CallbacksForKey_Errors(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.CallbacksForKey(Key(0))
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}

func TestCallbacks_CallbacksForKey_Success(t *testing.T) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		t.Fatalf("expected non-nil callbacks, got nil")
	}

	_, err := cbs.AddSingleShot(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	_, err = cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	cbSlice, err := cbs.CallbacksForKey(Key(0))
	if err != nil {
		t.Fatalf("expected nil error, got %s", err)
	}
	if len(cbSlice) != 2 {
		t.Fatalf("expected 2 callbacks, got %d", len(cbSlice))
	}
}

func BenchmarkContinuousCallbacks(b *testing.B) {
	cbs := New("Test", nil, nil)
	if cbs == nil {
		b.Fatalf("expected non-nil callbacks, got nil")
	}

	tag, err := cbs.AddContinuous(Key(0), func() {})
	if err != nil {
		b.Fatalf("expected nil error, got %q", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		cb, err := cbs.Callback(Key(0), tag)
		if err != nil {
			panic(err)
		}

		cb.(func())()
	}
}
