module github.com/karino2/folang/fc

go 1.23.4

replace github.com/karino2/folang/pkg/slice => ../pkg/slice

replace github.com/karino2/folang/pkg/buf => ../pkg/buf

replace github.com/karino2/folang/pkg/frt => ../pkg/frt

require (
	github.com/karino2/folang/pkg/frt v0.0.0-20250130153521-e095d25781c4
	github.com/karino2/folang/pkg/slice v0.0.0-20250130153521-e095d25781c4
)

require golang.org/x/exp v0.0.0-20250128182459-e0ece0dbea4c // indirect

replace github.com/karino2/folang/pkg/strings => ../pkg/strings
