//go:build linux
// +build linux

package protocoltalkers

import (
	"testing"

	"github.com/BBeRsErKeRR/system-stats-monitor/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a mock implementation of the executor function.
type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Execute() (chan string, chan error) {
	args := m.Called()
	return args.Get(0).(chan string), args.Get(1).(chan error)
}

// Test_parseTCPDumpOut tests the parseTCPDumpOut function.
func Test_parseTCPDumpOut(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		expected interface{}
		err      error
	}{
		{
			name: "valid line",
			line: "IP 192.168.1.1.123 > 192.168.1.2.456: UDP, length 60",
			expected: storage.ProtocolTalkerItem{
				Protocol:  "UDP",
				SendBytes: 60,
			},
			err: nil,
		},
		{
			name:     "empty line",
			line:     "",
			expected: nil,
			err:      nil,
		},
		{
			name:     "invalid line with less than 5 fields",
			line:     "IP 192.168.1.1.123 >",
			expected: nil,
			err:      nil,
		},
		{
			name:     "invalid line with non-IP protocol",
			line:     "ARP, Request who-has 10.0.0.1.42545 (00:25:4d:0c:e7:b1) tell 10.0.0.2.39510, length 28",
			expected: nil,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseTCPDumpOut(tc.line)
			if tc.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, actual)
		})
	}
}
