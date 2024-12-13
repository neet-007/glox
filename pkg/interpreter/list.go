package interpreter

import (
	"fmt"

	"github.com/neet-007/glox/pkg/parser"
	"github.com/neet-007/glox/pkg/runtime"
)

type List struct {
	list  parser.ListExpr
	items []any
}

func NewList(list parser.ListExpr, items []any) List {
	return List{
		list:  list,
		items: items,
	}
}

func (l List) Get(i int) (any, *runtime.RuntimeError) {
	if i >= len(l.items) {
		return nil, runtime.NewRuntimeError(l.list.LeftBracket, fmt.Sprintf("index out of bound index %d length %d", i, len(l.items)))
	}

	return l.items[i], nil
}

func (l List) Set(i int, value any) *runtime.RuntimeError {
	if i >= len(l.items) {
		return runtime.NewRuntimeError(l.list.LeftBracket, fmt.Sprintf("index out of bound index %d length %d", i, len(l.items)))
	}

	l.items[i] = value

	return nil
}

func (l List) Append(value any) {
	l.items = append(l.items, value)
}

func (l List) String() string {
	return "<list>"
}
