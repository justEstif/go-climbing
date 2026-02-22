package components

import "context"

type ContextKey string

const IsSignedInKey ContextKey = "isSignedIn"

func GetIsSignedIn(ctx context.Context) bool {
	if v, ok := ctx.Value(IsSignedInKey).(bool); ok {
		return v
	}
	return false
}
