#!/bin/sh

PKG_INFO=../pkg/pkg_all.foi

./tinyfo $PKG_INFO $1
go fmt gen_*.go
