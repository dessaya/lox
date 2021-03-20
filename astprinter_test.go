// +build ignore

package lox

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAstPrinter(t *testing.T) {
	expression := NewBinary(
		NewUnary(
			NewToken(MINUS, "-", nil, 1),
			NewLiteral(123),
		),
		NewToken(STAR, "*", nil, 1),
		NewGrouping(
			NewLiteral(45.67),
		),
	)

	require.Equal(t, "(* (- 123) (group 45.67))", ExprToString(expression))
}
