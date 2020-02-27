package main

import (
	"testing"
	"time"
)

func TestFileName(t *testing.T) {
	s := FileName(time.Now())
	t.Logf("%s", s)
}
