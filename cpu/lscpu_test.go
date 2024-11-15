package cpu

import (
	"encoding/json"
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestGetHostLsCpu(t *testing.T) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(hostLsCpu))
}

func TestParseHostLsCpu(t *testing.T) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		t.Fatalf(err.Error())
	}

	cpuInfo, err := parseLsCpu(hostLsCpu)
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, err := json.MarshalIndent(cpuInfo, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(jsonData))
}

var testFiles = []string{
	"test_data/dell-r430-lscpu.json",
	"test_data/hp-dl380p-gen8-lscpu.json",
	"test_data/rpi5-lscpu.json",
	"test_data/mediatek-genio-1200-lscpu.json",
	"test_data/mediatek-g350-lscpu.json",
	"test_data/intel-cbrd-raptor-lake.json",
	"test_data/intel-core2.json",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsCpuFile := range testFiles {
		t.Run(lsCpuFile, func(t *testing.T) {
			lsCpu, err := os.ReadFile(lsCpuFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			cpuInfo, err := parseLsCpu(lsCpu)
			if err != nil {
				t.Fatalf(err.Error())
			}

			jsonData, err := json.MarshalIndent(cpuInfo, "", "  ")
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log(string(jsonData))
		})
	}
}

// TestUtsName tests the Uname syscall to see what format the architecture is reported in
func TestUtsName(t *testing.T) {
	var sysInfo unix.Utsname
	err := unix.Uname(&sysInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(sysInfo.Sysname[:]))    // Linux
	t.Log(string(sysInfo.Nodename[:]))   // jpmeijers-XP-13-7390
	t.Log(string(sysInfo.Release[:]))    // 6.8.0-48-generic
	t.Log(string(sysInfo.Version[:]))    // #48-Ubuntu SMP PREEMPT_DYNAMIC Fri Sep 27 14:04:52 UTC 2024
	t.Log(string(sysInfo.Machine[:]))    // x86_64
	t.Log(string(sysInfo.Domainname[:])) // (none)
}
