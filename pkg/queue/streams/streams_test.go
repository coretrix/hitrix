package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xorcare/pointer"
)

func TestGetGroupName(t *testing.T) {
	type testCase struct {
		streamName        string
		suffix            *string
		expectedGroupName string
	}

	testCases := []testCase{
		{
			streamName:        "first",
			suffix:            nil,
			expectedGroupName: "first_group",
		}, {
			streamName:        "first",
			suffix:            pointer.String("name"),
			expectedGroupName: "first_group_name",
		},
	}

	for _, oneTestCase := range testCases {
		have := GetGroupName(oneTestCase.streamName, oneTestCase.suffix)
		assert.Equal(t, oneTestCase.expectedGroupName, have)
	}
}
