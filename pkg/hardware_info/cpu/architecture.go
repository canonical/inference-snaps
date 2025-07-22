package cpu

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/canonical/stack-utils/pkg/constants"
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
		"aarch64": constants.Arm64,
		"armv7l":  constants.Armhf,
		"armv8l":  constants.Arm64,
		"i686":    constants.I386,
		"ppc":     constants.Powerpc,
		"ppc64":   constants.Ppc64,
		"ppc64le": constants.Ppc64el,
		"riscv64": constants.Riscv64,
		"s390x":   constants.S390x,
		"x86_64":  constants.Amd64,
	}

	if debArch, ok := lookupTable[unameArch]; !ok {
		return "", fmt.Errorf("unsupported architecture: %s", unameArch)
	} else {
		return debArch, nil
	}

}
