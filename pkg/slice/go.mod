module github.com/karino2/folang/pkg/slice

go 1.23.4

require github.com/karino2/folang/pkg/frt v0.0.0-20250202140944-2d44ccfef24a

require (
	github.com/google/go-cmp v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c // indirect
)

replace github.com/karino2/folang => ../../

replace github.com/karino2/folang/pkg/frt => ../frt
