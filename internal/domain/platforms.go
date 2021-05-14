package domain

// Used for routing in usecases
const (
	AireOSPlatformName      string = "aireos"
	AlliedWarePlatformName  string = "allied"
	ASASoftwarePlatformName string = "asa"
	IOSPlatformName         string = "ios"
	IronWarePlatformName    string = "foundry"
	JunOSPlatformName       string = "srx"
	NXOSPlatformName        string = "nxos"
	ScreenOSPlatformName    string = "ssg"
	YamahaOSPlatformName    string = "yamaha"
)

// Used for config validation
var (
	Platforms = []string{
		AireOSPlatformName,
		AlliedWarePlatformName,
		ASASoftwarePlatformName,
		IOSPlatformName,
		IronWarePlatformName,
		JunOSPlatformName,
		NXOSPlatformName,
		ScreenOSPlatformName,
		YamahaOSPlatformName,
	}
)
