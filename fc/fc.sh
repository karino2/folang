#!/bin/sh

./tinyfo $1
go fmt gen_*.go
