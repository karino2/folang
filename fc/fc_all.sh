#!/bin/sh

PKG_INFO=../pkg/pkg_all.foi

./tinyfo $PKG_INFO ftype.fo ast.fo expr_to_type.fo expr_to_go.fo stmt_to_go.fo tokenizer.fo
go fmt gen_*.go
