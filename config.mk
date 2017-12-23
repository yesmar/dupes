cmd ?= dupes

version ?= 1.0.1

commit ?= $(shell git rev-parse --short HEAD 2>/dev/null)
commit := $(if $(commit),+$(commit),)

sources?=main.go
