package domain

// Used for routing in usecases
const (
	IOSPlatformName         string = "ios"
	IronWarePlatformName    string = "foundry"
	AireOSPlatformName      string = "aireos"
	AlliedWarePlatformName  string = "allied"
	ASASoftwarePlatformName string = "asa"
	ScreenOSPlatformName    string = "ssg"
	JunOSPlatformName       string = "srx"
	YamahaOSPlatformName    string = "yamaha"
)

// Used for config validation
var (
	Platforms = []string{
		IOSPlatformName,
		IronWarePlatformName,
		AireOSPlatformName,
		AlliedWarePlatformName,
		ASASoftwarePlatformName,
		ScreenOSPlatformName,
		YamahaOSPlatformName,
	}
)
