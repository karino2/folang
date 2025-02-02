#!/bin/sh

PKG_INFO=../pkg/pkg_all.foi

./tinyfo $PKG_INFO ftype.fo ast.fo
go fmt gen_*.go
