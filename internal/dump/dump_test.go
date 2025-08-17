package dump_test

import (
	"testing"

	"github.com/trentwiles/hackernews/internal/dump"
)

func TestCleanExport(t *testing.T) {
	dump.WipeExports()
}