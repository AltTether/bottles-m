package bottles

import (
	"testing"
	"context"

	//"github.com/stretchr/testify/assert"
)

func TestSetHandlersToEngine(t *testing.T) {
	engine := New()

	bottleAddHandlerFunc := func(ctx context.Context, b *Bottle) {}
	bottleGetHandlerFunc := func(ctx context.Context, b *Bottle) {}

	engine.SetBottleAddHandler(bottleAddHandlerFunc)
	engine.SetBottleGetHandler(bottleGetHandlerFunc)
}

func TestRunEngineWithNoHandler(t *testing.T) {
	engine := New()

	engine.Run()
	engine.Stop()
}

func TestAddFromGateway(t *testing.T) {
	engine := New()
	gateway := engine.Gateway

	engine.Run()

	gateway.Add(&Bottle{})

	engine.Stop()
}
