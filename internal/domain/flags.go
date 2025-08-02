package domain

// Flag names
const (
	HostnameFlagName        string = "hostname"
	PortFlagName            string = "port"
	TimeoutFlagName         string = "timeout"
	CommandFlagName         string = "command"
	ExecPlatformFlagName    string = "exec-platform"
	EnableModeFlagName      string = "enable-mode"
	RedundantModeFlagName   string = "redundant-mode"
	SecureModeFlagName      string = "secure-mode"
	DefaultPrivModeFlagName string = "default-privilege-mode"
	UsernameFlagName        string = "username"
	PasswordFlagName        string = "password"
	PrivPasswordFlagName    string = "priv-password"
	HostKeyPathFlagName     string = "host-key-path"
)

// Flag usages
const (
	HostnameFlagUsage        string = "Set hostname or IP address."
	PortFlagUsage            string = "Set port number."
	TimeoutFlagUsage         string = "Set timeout seconds."
	CommandFlagUsage         string = "Set a command."
	ExecPlatformFlagUsage    string = "Set exec-platform. Refer to README.md what to be set."
	EnableModeFlagUsage      string = "Raise to privileged EXEC mode."
	RedundantModeFlagUsage   string = "Use redundant prompt mode."
	SecureModeFlagUsage      string = "Use ssh mode."
	DefaultPrivModeFlagUsage string = "Use default privileged mode assinged by RADIUS attribute."
	UsernameFlagUsage        string = "Set username."
	PasswordFlagUsage        string = "Set password."
	PrivPasswordFlagUsage    string = "Set password to raise to privileged EXEC mode."
	HostKeyPathFlagUsage     string = "Set path to host key file for SSH host key verification."
)

// Flag values
const (
	PortFlagDefaultValue            int64  = 0
	TimeoutFlagDefaultValue         int64  = 5
	ExecPlatformFlagDefaultValue    string = "ios"
	EnableModeFlagDefaultValue      bool   = false
	RedundantModeFlagDefaultValue   bool   = false
	SecureModeFlagDefaultValue      bool   = false
	DefaultPrivModeFlagDefaultValue bool   = false
	UsernameFlagDefaultValue        string = "admin"
	PasswordFlagDefaultValue        string = "cisco"
	PrivPasswordFlagDefaultValue    string = "enable"
	HostKeyPathFlagDefaultValue     string = ""
)

// Flag requireds
const (
	HostnameFlagRequired bool = true
	CommandFlagRequired  bool = true
)

// Flag aliases
var (
	HostnameFlagAliases        = []string{"H"}
	PortFlagAliases            = []string{"P"}
	TimeoutFlagAliases         = []string{"t"}
	CommandFlagAliases         = []string{"C"}
	ExecPlatformFlagAliases    = []string{"x"}
	EnableModeFlagAliases      = []string{"e", "ena", "enable"}
	RedundantModeFlagAliases   = []string{"r", "redundant"}
	DefaultPrivModeFlagAliases = []string{"d"}
	SecureModeFlagAliases      = []string{"s", "sec", "secure"}
	UsernameFlagAliases        = []string{"u"}
	PasswordFlagAliases        = []string{"p"}
	PrivPasswordFlagAliases    = []string{"pp"}
	HostKeyPathFlagAliases     = []string{"hkp"}
)

// Flag envvars
var (
	HostnameFlagEnvVars     = "TELEE_HOSTNAME"
	CommandFlagEnvVars      = "TELEE_COMMAND"
	UsernameFlagEnvVars     = "TELEE_USERNAME"
	PasswordFlagEnvVars     = "TELEE_PASSWORD"
	PrivPasswordFlagEnvVars = "TELEE_PRIVPASSWORD"
	HostKeyPathFlagEnvVars  = "TELEE_HOSTKEYPATH"
)
