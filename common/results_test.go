package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_empty_results_cannot_be_summarized(t *testing.T) {
	_, err := CalculateResults(&mockFilesystem{}, &mockValueSource{})
	assert.Error(t, err)
}

func Test_duplicates_have_correct_ids(t *testing.T) {
	mockFs := new(mockFilesystem)

	// setup expectations
	mockFs.On("GetAbsolutePath", mock.Anything).Return("a")
	mockFs.On("PathStatus", mock.Anything).Return(true, true)

	rows, err := CalculateResults(mockFs, &mockValueSource{values_: []string{
		"a",
		"a",
	}})
	require.NoError(t, err)
	require.Len(t, rows, 2)

	expectedRow := ResultRow{
		Id:           1,
		Path:         "a",
		ExpandedPath: "a",
		Exists:       true,
		IsDir:        true,
		Duplicates:   []int{2},
	}
	assert.Equal(t, expectedRow, rows[0])

	// swapped id & duplicate ids expected
	expectedRow.Id = 2
	expectedRow.Duplicates = []int{1}
	assert.Equal(t, expectedRow, rows[1])
}

func Test_no_duplicates(t *testing.T) {
	mockFs := new(mockFilesystem)

	// setup expectations
	mockFs.On("GetAbsolutePath", "a").Return("a")
	mockFs.On("GetAbsolutePath", "b").Return("b")
	mockFs.On("PathStatus", mock.Anything).Return(true, true)

	rows, _ := CalculateResults(mockFs, &mockValueSource{values_: []string{
		"a",
		"b",
	}})
	assert.Empty(t, rows[0].Duplicates)
	assert.Empty(t, rows[1].Duplicates)
	mockFs.AssertExpectations(t)
}

type mockFilesystem struct {
	mock.Mock
}

func (m *mockFilesystem) GetAbsolutePath(path string) string {
	args := m.Called(path)
	return args.String(0)
}

func (m *mockFilesystem) PathStatus(path string) (bool, bool) {
	args := m.Called(path)
	return args.Bool(0), args.Bool(1)
}

type mockValueSource struct {
	source_ string
	orig_   string
	values_ []string
}

func (m *mockValueSource) Source() string {
	return m.source_
}

func (m *mockValueSource) Orig() string {
	return m.orig_
}

func (m *mockValueSource) Values() []string {
	return m.values_
}
