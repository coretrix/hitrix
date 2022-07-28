package depth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

func GetOperationDepth(ctx context.Context) int {
	return findSelectionDepth(graphql.GetOperationContext(ctx).Operation.SelectionSet)
}
func findSelectionDepth(selections ast.SelectionSet) int {
	maxDepth := 0

	for _, selection := range selections {
		if field, isField := selection.(*ast.Field); isField && field != nil {
			if len(field.SelectionSet) > 0 {
				if depth := findSelectionDepth(field.SelectionSet); depth+1 > maxDepth {
					maxDepth = depth + 1
				}
			}
		} else if fragment, isFragmentSpread := selection.(*ast.FragmentSpread); isFragmentSpread && fragment != nil {
			if len(fragment.Definition.SelectionSet) > 0 {
				if depth := findSelectionDepth(fragment.Definition.SelectionSet); depth+1 > maxDepth {
					maxDepth = depth + 1
				}
			}
		} else if inlineFragment, isInlineFragment := selection.(*ast.InlineFragment); isInlineFragment && inlineFragment != nil {
			if len(inlineFragment.SelectionSet) > 0 {
				if depth := findSelectionDepth(inlineFragment.SelectionSet); depth+1 > maxDepth {
					maxDepth = depth + 1
				}
			}
		}
	}

	return maxDepth
}
