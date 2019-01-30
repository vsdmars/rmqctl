package pkg

import "fmt"

type version struct {
	major int
	minor int
	patch int
}

func (v version) string() string {

	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

var rmqctlVersion = version{
	major: 1,
	minor: 0,
	patch: 7,
}
