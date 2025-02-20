#!/bin/sh

PKG_INFO=../../pkg/pkg_all.foi

./fc $PKG_INFO $1
go fmt gen_*.go
