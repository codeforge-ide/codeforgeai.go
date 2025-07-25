module github.com/codeforge-ide/codeforgeai.go

go 1.23.0

toolchain go1.24.4

require (
	github.com/codeforgeai/codeforgeai.go v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.9.1
)

require golang.org/x/sys v0.33.0 // indirect

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/term v0.32.0
)

replace github.com/codeforgeai/codeforgeai.go => ./
