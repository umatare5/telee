package domain

// Used for routing in usecases
const (
	IOSPlatformName        string = "ios"
	AireOSPlatformName     string = "aireos"
	AlliedWarePlatformName string = "allied"
	ScreenOSPlatformName   string = "ssg"
)

// Used for config validation
var (
	CmdPlatforms = []string{IOSPlatformName, AireOSPlatformName, AlliedWarePlatformName, ScreenOSPlatformName}
)
