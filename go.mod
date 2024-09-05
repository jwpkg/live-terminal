module github.com/jwpkg/living-terminal


go 1.23.0

replace github.com/jwpkg/living-terminal/components => ./components

require golang.org/x/term v0.23.0

require (
	github.com/gabe565/go-spinners v1.1.0 // indirect
	github.com/muesli/cancelreader v0.2.2
	golang.org/x/sys v0.24.0 // indirect
)
