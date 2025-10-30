package dto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapQueryToListLogsFilter(t *testing.T) {
	testcases := []struct {
		name        string
		queries     map[string]string
		expectedErr error
		out         *ListLogsFilter
	}{
		{
			name: "ok",
			queries: map[string]string{
				"limit":      "1",
				"page":       "1",
				"module_id":  "1",
				"exploit_id": "a207bb9c-3bf3-4b57-b782-a6aba66da239",
			},
			out: &ListLogsFilter{
				Limit:     1,
				Page:      1,
				ModuleId:  1,
				ExploitId: "a207bb9c-3bf3-4b57-b782-a6aba66da239",
			},
		},
		{
			name: "negative page",
			queries: map[string]string{
				"limit":      "1",
				"page":       "-1",
				"module_id":  "1",
				"exploit_id": "a207bb9c-3bf3-4b57-b782-a6aba66da239",
			},
			expectedErr: ErrPageNegative,
		},
		{
			name: "incorrect page",
			queries: map[string]string{
				"limit":      "1",
				"page":       "-dasdsa1",
				"module_id":  "1",
				"exploit_id": "a207bb9c-3bf3-4b57-b782-a6aba66da239",
			},
			expectedErr: ErrPageIncorrect,
		},
		{
			name: "incorrect exploit_id",
			queries: map[string]string{
				"limit":      "1",
				"page":       "1",
				"module_id":  "1",
				"exploit_id": "1das",
			},
			expectedErr: ErrExploitIdIncorrect,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(tt *testing.T) {
			out, err := MapQueryToListLogsFilter(tc.queries)
			require.ErrorIs(tt, err, tc.expectedErr)
			require.Equal(tt, out, tc.out)
		})
	}
}
