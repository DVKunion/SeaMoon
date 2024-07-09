package version

import core "github.com/v2fly/v2ray-core/v5"

var (
	Version          string = "dev"
	Commit           string = ""
	V2rayCoreVersion string = core.Version()
)
