cmd := dupes

version := 1.1.0
commit := $(shell git rev-parse --short HEAD 2>/dev/null)
release := $(if $(commit),$(version)+$(commit),$(version))

sources?=main.go
