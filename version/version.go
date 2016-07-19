package version

// Version identifies this Deis product revision.
const Version = "2.2.0"

// BuildVersion is the git revision of the build.
// Note: This value is overwritten by the linker during build
var BuildVersion = ""
