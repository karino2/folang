module github.com/karino2/folang/samples

go 1.23.4

replace github.com/karino2/folang/pkg/frt => ../pkg/frt

replace github.com/karino2/folang/pkg/buf => ../pkg/buf

require (
	github.com/karino2/folang/pkg/buf v0.0.0-00010101000000-000000000000
	github.com/karino2/folang/pkg/frt v0.0.0-00010101000000-000000000000
	github.com/karino2/folang/pkg/slice v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c // indirect
)

replace github.com/karino2/folang/pkg/slice => ../pkg/slice

replace github.com/karino2/folang/pkg/strings => ../pkg/strings
