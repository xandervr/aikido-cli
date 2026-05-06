package version

// Version, Commit, and Date are set at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

type Info struct {
	Version string `json:"version" aikido:"column,header=Version"`
	Commit  string `json:"commit"  aikido:"column,header=Commit"`
	Date    string `json:"date"    aikido:"column,header=Date"`
}

func Current() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}
