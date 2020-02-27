package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/dyweb/gommon/util/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFileName(t *testing.T) {
	s := FileName(time.Now())
	t.Logf("%s", s)
}

func TestStripHeader(t *testing.T) {
	b := testutil.ReadFixture(t, "testdata/2020-01-06-weekly.md")
	s := StripHeader(b)
	assert.False(t, bytes.Contains(s, frontMatterDash))
}
