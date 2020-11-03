package version

var (
	// Version is an app's version, filled from ldflags at compile
	Version string = ""
	// BuildTime is a compilation's datetime of an app
	BuildTime string = ""
)
