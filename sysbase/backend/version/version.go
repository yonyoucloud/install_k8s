package version

import (
	"fmt"
)

var (
	Pkg           = ""
	Version       = ""
	GitCommitSha  = ""
	GitCommitTime = ""
	BuildTime     = ""
)

func VersionInfo() string {
	return fmt.Sprintf("Sysbase git:%s, version:%s, commit-sha:%s, commit-time:%s, build-time:%s", Pkg, Version, GitCommitSha, GitCommitTime, BuildTime)
}
