package helpers

import (
	"runtime"
)

const ColorRed string = "\033[31m"
const ColorGreen string = "\033[32m"
const ColorCyan string = "\033[34m"
const ColorYellow string = "\033[33m"
const ColorReset string = "\033[0m"

func CurrentPlatform() string {
	switch runtime.GOOS {
	case "darwin":
		// macOS
		return "macos"
	case "linux":
		if Exists("/etc/lsb-release") {
			return "ubuntu"
		} else if Exists("/etc/debian_version") {
			return "debian"
		} else if Exists("/etc/redhat-release") {
			return "redhat"
		} else if Exists("/etc/centos-release") {
			return "centos"
		} else if Exists("/etc/fedora-release") {
			return "fedora"
		}
	case "freebsd":
		return "freebsd"
	}
	return "unsupported"
}

func DefaultPkgManager() string {
	switch runtime.GOOS {
	case "darwin":
		// macOS
		return "homebrew"
	case "linux":
		if Exists("/etc/lsb-release") || Exists("/etc/debian_version") {
			// Ubuntu, Mint, Debian
			return "apt-get"
		} else if Exists("/etc/redhat-release") || Exists("/etc/centos-release") || Exists("/etc/fedora-release") {
			// Red Hat, CentOS, Fedora
			return "rpm"
		}
	case "freebsd":
		return "pkg"
	}
	return "unsupported"
}
