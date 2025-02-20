module github.com/karino2/folang/cmd/build_sample_md

go 1.23.4

replace github.com/karino2/folang/pkg/dict => ../../pkg/dict

replace github.com/karino2/folang/pkg/frt => ../../pkg/frt

replace github.com/karino2/folang/pkg/strings => ../../pkg/strings

replace github.com/karino2/folang/pkg/slice => ../../pkg/slice

replace github.com/karino2/folang/pkg/sys => ../../pkg/sys

replace github.com/karino2/folang/pkg/buf => ../../pkg/buf

require (
	github.com/karino2/folang/pkg/buf v0.0.0-20250220075419-c2c06e428e77
	github.com/karino2/folang/pkg/frt v0.0.0-20250220075419-c2c06e428e77
	github.com/karino2/folang/pkg/slice v0.0.0-20250220075419-c2c06e428e77
	github.com/karino2/folang/pkg/strings v0.0.0-20250220075419-c2c06e428e77
	github.com/karino2/folang/pkg/sys v0.0.0-20250220075419-c2c06e428e77
)

require github.com/google/go-cmp v0.6.0 // indirect

replace github.com/karino2/folang => ../..
