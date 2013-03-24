package goku

import (
    "testing"
    // "fmt"
    "github.com/sdegutis/go.assert"
)

func TestMemorySession(t *testing.T) {
    sid := "test-session"
    p := memoryProvider{}
    s, err := p.StartSession(sid)
    assert.Equals(t, err, nil)
    key := "test-key"
    value := "test-value"
    v := s.Get(key)
    assert.Equals(t, v, nil)

    err = s.Set(key, value)
    assert.Equals(t, err, nil)

    p2 := memoryProvider{}
    s, err = p2.StartSession(sid)
    assert.Equals(t, err, nil)

    v = s.Get(key)
    assert.NotEquals(t, v, nil)
    assert.Equals(t, v.(string), value)
}
