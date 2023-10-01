package jtdinfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHintSet(t *testing.T) {
	hs := HintSet{
		Values: [][]string{
			{"a", "b", "c"},
		},
	}

	assert.False(t, hs.IsActive())

	v, found := hs.PeekActive()
	assert.Empty(t, v)
	assert.False(t, found)

	hsA := hs.SubHints("a")
	assert.False(t, hsA.IsActive())

	v, found = hsA.PeekActive()
	assert.Empty(t, v)
	assert.False(t, found)

	hsB := hs.SubHints("a").SubHints("b")
	assert.False(t, hsB.IsActive())

	v, found = hsB.PeekActive()
	assert.Equal(t, "c", v)
	assert.True(t, found)

	hsC := hs.SubHints("a").SubHints("b").SubHints("c")
	assert.True(t, hsC.IsActive())

	v, found = hsC.PeekActive()
	assert.Empty(t, v)
	assert.False(t, found)
}

func TestHintSetWildcard(t *testing.T) {
	hs := HintSet{
		Values: [][]string{
			{"a", "b", "c"},
			{"d", "-", "e"},
		},
	}

	assert.False(t, hs.SubHints("a").SubHints("x").SubHints("c").IsActive())
	assert.True(t, hs.SubHints("d").SubHints("x").SubHints("e").IsActive())
}
