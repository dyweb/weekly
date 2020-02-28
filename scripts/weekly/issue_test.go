package main

import (
	"context"
	"testing"
)

func TestCloseOld(t *testing.T) {
	CloseOld(context.Background(), true)
}

func TestOpenNew(t *testing.T) {
	OpenNew(context.Background(), true)
}
