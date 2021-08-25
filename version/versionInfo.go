package version

import "fmt"

var (
	Name    string
	Version string
	Commit  string
	Built   string
)

func BuildInfo() []string {
	return []string{
		fmt.Sprintf("Name     : %s", Name),
		fmt.Sprintf("Version  : %s", Version),
		fmt.Sprintf("Commit   : %s", Commit),
		fmt.Sprintf("Built    : %s", Built),
	}
}
