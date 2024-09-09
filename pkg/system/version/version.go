package version

import (
	"github.com/DVKunion/SeaMoon/plugins/xray"
)

var (
	Version     string = "dev"
	Commit      string = ""
	XrayVersion string = xray.GetVer()
)
