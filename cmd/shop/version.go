package main

import (
	"fmt"
	"runtime/debug"
)

// version is set at release time by goreleaser via -ldflags.
var version = "dev"

// versionString returns the release version, or the module version when
// installed with "go install".
func versionString() string {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return version
}

func printVersion() {
	fmt.Println("shop " + versionString())
}
