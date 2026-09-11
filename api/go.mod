module github.com/sudzekai/web-os-api

go 1.26.8

require (
    github.com/sudzekai/web-os-api/logging v0.0.0 // direct
    github.com/sudzekai/web-os-api/server v0.0.0 // direct
    github.com/sudzekai/web-os-api/module v0.0.0 // direct
)

replace github.com/sudzekai/web-os-api/logging => ./packages/logging
replace github.com/sudzekai/web-os-api/server => ./packages/server
replace github.com/sudzekai/web-os-api/module => ./packages/module
