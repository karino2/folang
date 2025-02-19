#!/bin/sh

PKG_INFO=../pkg/pkg_all.foi

./fc $PKG_INFO ftype.fo ast.fo expr_to_type.fo expr_to_go.fo stmt_to_go.fo tokenizer.fo parse_state.fo infer.fo parser.fo main.fo
go fmt gen_*.go
