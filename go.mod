module github.com/cdzombak/libwx

go 1.22

toolchain go1.23.4

replace gonum.org/v1/gonum => github.com/cdzombak/gonum v0.0.0-20250111220929-90b1f7503766

require (
	github.com/stretchr/testify v1.9.0
	gonum.org/v1/gonum v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
