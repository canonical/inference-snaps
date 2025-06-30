package cpu

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	arm64   = "arm64"
	amd64   = "amd64"
	armhf   = "armhf"
	i386    = "i386"
	powerpc = "powerpc"
	ppc64   = "ppc64"
	ppc64el = "ppc64el"
	riscv64 = "riscv64"
	s390x   = "s390x"
)

func hostUnameMachine() (string, error) {
	// uname --machine
	out, err := exec.Command("uname", "--machine").Output()
	if err != nil {
		return "", err
	}
	architecture := string(out)
	return strings.TrimSpace(architecture), nil
}

// debianArchitecture translates the kernel architecture as reported by uname() to the debian architecture
// Based on lookup table from snapd: https://github.com/canonical/snapd/blob/master/arch/arch.go
func debianArchitecture(unameArch string) (string, error) {
	// Trim whitespace
	unameArch = strings.TrimSpace(unameArch)

	lookupTable := map[string]string{
		// uname:  debian
		"aarch64": arm64,
		"armv7l":  armhf,
		"armv8l":  arm64,
		"i686":    i386,
		"ppc":     powerpc,
		"ppc64":   ppc64,
		"ppc64le": ppc64el,
		"riscv64": riscv64,
		"s390x":   s390x,
		"x86_64":  amd64,
	}

	if debArch, ok := lookupTable[unameArch]; !ok {
		return "", fmt.Errorf("unsupported architecture: %s", unameArch)
	} else {
		return debArch, nil
	}

}
