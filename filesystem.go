package main

type filesystem interface {
	getAbsolutePath(path string) string
	pathStatus(path string) (bool, bool)
}
