package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_empty_results_cannot_be_summarized(t *testing.T) {
	_, err := calculateResults(&mockFilesystem{}, &mockValueSource{})
	assert.Error(t, err)
}

func Test_duplicates_have_correct_ids(t *testing.T) {
	mockFs := new(mockFilesystem)

	// setup expectations
	mockFs.On("getAbsolutePath", mock.Anything).Return("a")
	mockFs.On("pathStatus", mock.Anything).Return(true, true)

	rows, err := calculateResults(mockFs, &mockValueSource{values_: []string{
		"a",
		"a",
	}})
	require.NoError(t, err)
	require.Len(t, rows, 2)

	expectedRow := resultRow{
		id:           1,
		path:         "a",
		expandedPath: "a",
		exists:       true,
		isDir:        true,
		duplicates:   []int{2},
	}
	assert.Equal(t, expectedRow, rows[0])

	// swapped id & duplicate ids expected
	expectedRow.id = 2
	expectedRow.duplicates = []int{1}
	assert.Equal(t, expectedRow, rows[1])
}

func Test_no_duplicates(t *testing.T) {
	mockFs := new(mockFilesystem)

	// setup expectations
	mockFs.On("getAbsolutePath", "a").Return("a")
	mockFs.On("getAbsolutePath", "b").Return("b")
	mockFs.On("pathStatus", mock.Anything).Return(true, true)

	rows, _ := calculateResults(mockFs, &mockValueSource{values_: []string{
		"a",
		"b",
	}})
	assert.Empty(t, rows[0].duplicates)
	assert.Empty(t, rows[1].duplicates)
	mockFs.AssertExpectations(t)
}

type mockFilesystem struct {
	mock.Mock
}

func (m *mockFilesystem) getAbsolutePath(path string) string {
	args := m.Called(path)
	return args.String(0)
}

func (m *mockFilesystem) pathStatus(path string) (bool, bool) {
	args := m.Called(path)
	return args.Bool(0), args.Bool(1)
}

type mockValueSource struct {
	source_ string
	orig_   string
	values_ []string
}

func (m *mockValueSource) source() string {
	return m.source_
}

func (m *mockValueSource) orig() string {
	return m.orig_
}

func (m *mockValueSource) values() []string {
	return m.values_
}
