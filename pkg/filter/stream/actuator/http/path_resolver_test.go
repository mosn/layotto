package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pathResolver_isEmpty(t *testing.T) {
	p := NewPathResolver("/")
	assert.True(t, !p.HasNext())
	p = NewPathResolver("")
	assert.True(t, !p.HasNext())
	p = NewPathResolver("/a")
	assert.True(t, p.HasNext())
}

func Test_pathResolver_next(t *testing.T) {
	p := NewPathResolver("/a/b/c")
	assert.True(t, p.Next() == "a")
	assert.True(t, p.Next() == "b")
	assert.True(t, p.Next() == "c")
	assert.True(t, !p.HasNext())
	assert.True(t, p.Next() == "")
}

func Test_pathResolver_unresolvedPath(t *testing.T) {
	p := NewPathResolver("/a/b/c")
	assert.True(t, p.Next() == "a")
	assert.True(t, p.UnresolvedPath() == "/b/c")
	assert.True(t, p.Next() == "b")
	assert.True(t, p.UnresolvedPath() == "/c")
	assert.True(t, p.Next() == "c")
	assert.True(t, p.UnresolvedPath() == "")
	assert.True(t, !p.HasNext())
	assert.True(t, p.Next() == "")
	assert.True(t, p.UnresolvedPath() == "")
}
