package config

import (
	"path"
	"strings"
	"time"
)

// Compile time constants that should not be configurable
// during runtime.
const (
	Name        = "usersms"
	VersionFull = "0.1.0" // Use http://semver.org standards
	Description = "User profile management micro-service"

	// be sure to run micro with the new namespace after changing this e.g.
	//     $micro api --namespace=new.namespace.value ...
	// or set the environment value.
	// Docs here: https://micro.mu/docs/api.html#set-namespace
	Namespace = "go.micro.api"

	RPCNamePrefix = ""

	TimeFormat = time.RFC3339

	DocsPath = "docs"
)

var (
	// FIXME Probably won't work for none-unix systems!
	defaultInstallDir       = path.Join("/usr", "local", "bin")
	defaultSysDUnitFilePath = path.Join("/etc", "systemd", "system", DefaultSysDUnitName())
	sysDConfDir             = path.Join("/etc", Name)
	defaultConfDir          = sysDConfDir
)

func CanonicalName() string {
	return Name + VersionMajorPrefixed()
}

func CanonicalRPCName() string {
	return RPCNamePrefix + CanonicalName()
}

func VersionMajorPrefixed() string {
	return "v" + strings.SplitN(VersionFull, ".", 2)[0]
}

func WebNamePrefix() string {
	return Namespace + "." + VersionMajorPrefixed() + "."
}

func WebRootPath() string {
	return "/" + VersionMajorPrefixed() + "/" + Name
}

func CanonicalWebName() string {
	return WebNamePrefix() + Name
}

func DefaultSysDUnitName() string {
	return CanonicalName() + ".service"
}

func DefaultInstallDir() string {
	return defaultInstallDir
}

func DefaultInstallPath() string {
	return path.Join(defaultInstallDir, CanonicalName())
}

func DefaultSysDUnitFilePath() string {
	return defaultSysDUnitFilePath
}

func SysDConfDir() string {
	return sysDConfDir
}

// DefaultConfDir sets the value of the conf dir to use and returns it.
// It falls back to default - sysDConfDir - if newPSegments has zero len.
func DefaultConfDir(newPSegments ...string) string {
	if len(newPSegments) == 0 {
		defaultConfDir = sysDConfDir
	} else {
		defaultConfDir = path.Join(newPSegments...)
	}
	return defaultConfDir
}

func DefaultDocsDir() string {
	return path.Join(defaultConfDir, DocsPath)
}

func DefaultConfPath() string {
	return path.Join(defaultConfDir, CanonicalName()+".conf.yml")
}
