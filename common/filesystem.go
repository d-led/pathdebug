package common

type Filesystem interface {
	GetAbsolutePath(path string) string
	PathStatus(path string) (bool, bool)
}
