package static

import "embed"

//go:embed public/*
var AssetFiles embed.FS

//go:embed templates/*
var HtmlFiles embed.FS
