package selector

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/stack-utils/pkg/hardware_info"
	"github.com/canonical/stack-utils/pkg/types"
	"gopkg.in/yaml.v3"
)

// Test that the expected stack is chosen from a list of stacks
// device, stacks, expected stack

type machineTopStack struct {
	machine  string
	stacks   []string
	topStack string
}

var topStackSets = []machineTopStack{
	{
		// Ampere One machine should use ampere stack, not generic arm neon stack
		machine:  "ampere-one-m-banshee-12",
		stacks:   []string{"ampere", "arm-neon"},
		topStack: "ampere",
	},
	{
		// Ampere Altra machine should use ampere-altra stack, not generic arm neon stack
		machine:  "hp-proliant-rl300-gen11-altra",
		stacks:   []string{"ampere-altra", "arm-neon"},
		topStack: "ampere-altra",
	},
	{
		// Old CPU with Intel dGPU and NVIDIA dGPU - using intel-gpu because nvidia requires newer CPU flags
		machine:  "i5-3570k+arc-a580+gtx1080ti",
		stacks:   []string{"cpu-avx1", "cuda-generic", "intel-cpu", "intel-gpu"},
		topStack: "intel-gpu",
	},
	{
		// Machine with Intel CPU and Intel GPU should use GPU
		machine:  "mustang",
		stacks:   []string{"cpu-avx1", "cpu-avx2", "cpu-avx512", "intel-cpu", "intel-gpu"},
		topStack: "intel-gpu",
	},
	{
		// Machine with Intel iGPU and NVIDIA dGPU - always try and offload to dGPU if possible
		machine:  "system76-addw4",
		stacks:   []string{"cpu-avx1", "cpu-avx2", "cuda-generic", "intel-cpu", "intel-gpu"},
		topStack: "cuda-generic",
	},
	{
		// Machine with avx2 should prefer avx2 stack
		machine:  "xps13-7390",
		stacks:   []string{"cpu-avx1", "cpu-avx2"},
		topStack: "cpu-avx2",
	},
	{
		// Machine with Intel CPU should prefer intel-cpu stack above generic cpu stacks
		machine:  "xps13-7390",
		stacks:   []string{"cpu-avx1", "cpu-avx2", "intel-cpu"},
		topStack: "intel-cpu",
	},
}

func TestTopStack(t *testing.T) {
	for _, testSet := range topStackSets {
		t.Run(testSet.machine+"/"+testSet.topStack, func(t *testing.T) {
			var stacks []types.Stack
			for _, stackName := range testSet.stacks {
				stackManifestFile := fmt.Sprintf("../../test_data/engines/%s/engine.yaml", stackName)
				data, err := os.ReadFile(stackManifestFile)
				if err != nil {
					t.Fatal(err)
				}

				var stack types.Stack
				err = yaml.Unmarshal(data, &stack)
				if err != nil {
					t.Fatal(err)
				}

				stacks = append(stacks, stack)
			}

			hardwareInfo, err := hardware_info.GetFromRawData(t, testSet.machine, true)
			if err != nil {
				t.Fatal(err)
			}

			scoredStacks, err := ScoreStacks(hardwareInfo, stacks)
			if err != nil {
				t.Fatal(err)
			}

			topStack, err := TopStack(scoredStacks)
			if err != nil {
				t.Fatal(err)
			}

			if topStack.Name != testSet.topStack {
				for _, stack := range scoredStacks {
					t.Logf("%s=%d %s", stack.Name, stack.Score, strings.Join(stack.Notes, ", "))
				}
				t.Errorf("Top stack name: %s, expected: %s", topStack.Name, testSet.topStack)
			}
		})
	}
}
