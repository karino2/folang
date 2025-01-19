#!/bin/sh

./folang $1
go fmt *gen.go
