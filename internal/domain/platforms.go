package domain

// Used for routing in usecases
const (
	IOSPlatformName         string = "ios"
	NXOSPlatformName        string = "nxos"
	IronWarePlatformName    string = "foundry"
	AireOSPlatformName      string = "aireos"
	AlliedWarePlatformName  string = "allied"
	ASASoftwarePlatformName string = "asa"
	JunOSPlatformName       string = "srx"
	ScreenOSPlatformName    string = "ssg"
	YamahaOSPlatformName    string = "yamaha"
)

// Used for config validation
var (
	Platforms = []string{
		IOSPlatformName,
		NXOSPlatformName,
		IronWarePlatformName,
		AireOSPlatformName,
		AlliedWarePlatformName,
		ASASoftwarePlatformName,
		JunOSPlatformName,
		ScreenOSPlatformName,
		YamahaOSPlatformName,
	}
)
