package domain

// Used for routing in usecases
const (
	IOSPlatformName           string = "ios"
	IronWarePlatformName      string = "foundry"
	AireOSPlatformName        string = "aireos"
	AlliedWarePlatformName    string = "allied"
	ASASoftwarePlatformName   string = "asa"
	ASASoftwareHAPlatformName string = "asa-ha"
	ScreenOSPlatformName      string = "ssg"
)

// Used for config validation
var (
	CmdPlatforms = []string{
		IOSPlatformName,
		IronWarePlatformName,
		AireOSPlatformName,
		AlliedWarePlatformName,
		ASASoftwarePlatformName,
		ASASoftwareHAPlatformName,
		ScreenOSPlatformName,
	}
)
