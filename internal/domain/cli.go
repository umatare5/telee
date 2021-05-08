package domain

// Used for routing in usecases
const (
	IOSPlatformName    string = "ios"
	AireOSPlatformName string = "aireos"
)

// Used for config validation
var (
	CmdPlatforms = []string{IOSPlatformName, AireOSPlatformName}
)
