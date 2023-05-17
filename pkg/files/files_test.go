package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadlines(t *testing.T) {
	ret, err := ReadLines("files_test.go")
	require.NoError(t, err)
	require.Contains(t, ret[0], "package files", "could not read correctly")
}

func TestReadLinesOffsetN(t *testing.T) {
	ret, err := ReadLinesOffsetN("files_test.go", 2, 1)
	require.NoError(t, err)
	require.Contains(t, ret[0], "import (", "could not read correctly")
}

func TestReadFile(t *testing.T) {
	ret, err := ReadFile("files_test.go")
	require.NoError(t, err)
	require.Contains(t, ret, "TestReadFile(", "could not read correctly")
}
