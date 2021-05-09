module github.com/arindas/elevate

go 1.16

require (
    github.com/arindas/elevate/internal/http v0.0.0
    github.com/arindas/elevate/internal/app v0.0.0
)

replace github.com/arindas/elevate/internal/http => ./internal/http
replace github.com/arindas/elevate/internal/app => ./internal/app
