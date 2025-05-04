package snippets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentPrefix(t *testing.T) {
	assert.Equal(t, "//", commentPrefix(".js"))
}
