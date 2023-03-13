package main

type valueSource interface {
	source() string
	orig() string
	values() []string
}
