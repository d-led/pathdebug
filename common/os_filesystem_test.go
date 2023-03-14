package common

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sut Filesystem = &OsFilesystem{}

func Test_relative_paths_are_expanded(t *testing.T) {
	assert.NotContains(t, "..", sut.GetAbsolutePath(".."))
}

const interpolatedVar = "INTERPOLATED_VAR"
const interpolatedValue = "interpolated_value"

func Test_environment_vars_are_expanded(t *testing.T) {
	setInterpolatedVar()
	part := "some_path"
	fullPath := path.Join(getEnvVarPath(), part)

	absPath := sut.GetAbsolutePath(fullPath)

	assert.Contains(t, absPath, part)
	assert.Contains(t, absPath, interpolatedValue)
	assert.NotContains(t, absPath, interpolatedVar)
}

func Test_nonexistent_paths(t *testing.T) {
	exists, _ := sut.PathStatus(sut.GetAbsolutePath("some_nonexistent_path"))
	assert.False(t, exists)
}

func Test_existing_directories(t *testing.T) {
	// current path is assumed to exist
	exists, isDir := sut.PathStatus(sut.GetAbsolutePath("."))

	assert.True(t, exists)
	assert.True(t, isDir)
}

func Test_existing_files(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	exists, isDir := sut.PathStatus(tempFile.Name())

	assert.True(t, exists)
	assert.False(t, isDir)
}

func setInterpolatedVar() {
	os.Setenv(interpolatedVar, interpolatedValue)
}

func getEnvVarPath() string {
	if runtime.GOOS == "windows" {
		return "%" + interpolatedVar + "%"
	}
	return fmt.Sprintf("${%s}", interpolatedVar)
}

func createTempFile(t *testing.T) *os.File {
	tempFile, err := os.CreateTemp("", "some_tmp_file")
	require.NoError(t, err)
	fmt.Println(tempFile.Name())
	return tempFile
}
