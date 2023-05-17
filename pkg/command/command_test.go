package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	ret, err := WithContext(ctx, "ls")
	require.NoError(t, err)
	require.Contains(t, string(ret), "command_test.go", "exec incorrect")
}

func TestBadWithContext(t *testing.T) {
	ctx := context.Background()
	_, err := WithContext(ctx, "not_found_ls")
	require.Error(t, err)
}

func TestStream(t *testing.T) {
	ctx := context.Background()
	chanO, chanErr := Stream(ctx, "ls", "-lah")
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-chanErr:
			require.NoError(t, err)
		case in, ok := <-chanO:
			if !ok {
				return
			}
			require.NotEmpty(t, in)
		}
	}
}

func TestBadStream(t *testing.T) {
	ctx := context.Background()
	chanO, chanErr := Stream(ctx, "not_found_ls")
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-chanErr:
			require.Error(t, err)
		case in, ok := <-chanO:
			if !ok {
				return
			}
			require.NotEmpty(t, in)
		}
	}
}
