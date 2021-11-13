module github.com/arindas/elevate

go 1.16

require (
    github.com/arindas/elevate/pkg/http v0.0.0
    github.com/arindas/elevate/pkg/app v0.0.0
)

replace github.com/arindas/elevate/pkg/http => ./pkg/http
replace github.com/arindas/elevate/pkg/app => ./pkg/app
