module github.com/sudzekai/web-os-api/modules/files-module

go 1.26.8

require github.com/sudzekai/web-os-api/logging v0.0.0 // direct

replace github.com/sudzekai/web-os-api/logging => ../../packages/logging

require github.com/sudzekai/web-os-api/server v0.0.0 // direct

replace github.com/sudzekai/web-os-api/server => ../../packages/server
