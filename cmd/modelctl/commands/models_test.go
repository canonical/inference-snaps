package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func prepareModelsTestData() (*modelsCommand, *outputModels, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveModel("4b-it-int4-fq-ov")
	if err != nil {
		return nil, nil, fmt.Errorf("error setting active model name: %v", err)
	}
	err = cache.SetActiveEngine("intel-gpu")
	if err != nil {
		return nil, nil, fmt.Errorf("error setting active engine name: %v", err)
	}

	ctx := &common.Context{
		ModelsDir:  "../../../test_data/models",
		EnginesDir: "../../../test_data/engines",
		Cache:      cache,
		Config:     nil,
	}

	allModels, err := common.GetAllModels(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting all models with engines: %v", err)
	}

	cmd := modelsCommand{Context: ctx}

	activeModel, err := cmd.Cache.GetActiveModel()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %v", common.LookingUpActiveModel, err)
	}

	modelsList := outputModels{
		ActiveModel: activeModel,
		Models:      allModels,
	}

	return &cmd, &modelsList, nil
}

func TestModelsJson(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		t.Fatal(err)
	}
}

func TestModelsTable(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printModelsTable(*modelsList)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetModelsTable(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	tableStr, err := cmd.getModelsTable(*modelsList)
	if err != nil {
		t.Fatalf("Error getting models table: %v", err)
	}

	expectedTable := `NAME                 CAPABILITIES               DISK 
26b-q4-k-m-gguf      text                       6G   
30b-a3b-q4-k-m-gguf  text, vision, audio, tool  6G   
4b-it-int4-fq-ov*    text                       6G   
`

	if tableStr != expectedTable {
		t.Errorf("Models table not as expected.\n\nGot:\n\n%s\n\nWant:\n\n%s", tableStr, expectedTable)
	}
}

func TestGetModelsTableAllModels(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}
	cmd.all = true

	tableStr, err := cmd.getModelsTable(*modelsList)
	if err != nil {
		t.Fatalf("Error getting models table: %v", err)
	}

	expectedTable := `NAME                 CAPABILITIES               DISK   ENGINES                                                          
26b-q4-k-m-gguf      text                       6G     arm-neon, cpu, cpu-avx1, cpu-avx2, cpu-avx512, cuda-generic, ro… 
30b-a3b-q4-k-m-gguf  text, vision, audio, tool  6G     cpu, cuda-generic, rocm-generic                                  
4b-it-int4-fq-ov*    text                       6G     intel-cpu, intel-gpu, intel-npu                                  
`

	if tableStr != expectedTable {
		t.Errorf("Models table not as expected.\n\nGot:\n\n%s\n\nWant:\n\n%s", tableStr, expectedTable)
	}
}

func TestGetModelsTableIncludesHintForIncompatibleModels(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}
	cmd.all = false
	modelsList.Models = modelsList.Models[:2]

	tableStr, err := cmd.getModelsTable(*modelsList)
	if err != nil {
		t.Fatalf("Error getting models table: %v", err)
	}
	if !strings.Contains(tableStr, "Hint:") {
		t.Fatalf("expected hint in table output, got:\n%s", tableStr)
	}
	if !strings.Contains(tableStr, "models --all") {
		t.Fatalf("expected --all hint in output, got:\n%s", tableStr)
	}
}

func Example_printModelsJson() {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		panic(fmt.Sprintf("Error preparing test data: %v", err))
	}

	// Use only the 4b-it-int4-fq-ov model to keep output concise
	var filtered []common.ModelDetails
	for _, m := range modelsList.Models {
		if m.Name == "4b-it-int4-fq-ov" {
			filtered = append(filtered, m)
		}
	}
	modelsList.Models = filtered

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		panic(fmt.Sprintf("Error printing models json: %v", err))
	}

	// Output:
	// {
	//   "active-model": "4b-it-int4-fq-ov",
	//   "models": [
	//     {
	//       "name": "4b-it-int4-fq-ov",
	//       "alias": "4b-it",
	//       "description": "OpenVino 4b test model",
	//       "model-card-url": "https://example.com/model-card",
	//       "format": "OpenVINO_IR",
	//       "quantization": "int4-fq",
	//       "capabilities": [
	//         "text"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-4b-it-int4-fq-ov"
	//       ],
	//       "compatible-engines": [
	//         "intel-cpu",
	//         "intel-gpu",
	//         "intel-npu"
	//       ]
	//     }
	//   ]
	// }
}

func Example_printAllModelsJson() {
	cmd, modelsList, err := prepareModelsTestData()
	cmd.all = true
	if err != nil {
		panic(fmt.Sprintf("Error preparing test data: %v", err))
	}

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		panic(fmt.Sprintf("Error printing models json: %v", err))
	}

	// Output:
	// {
	//   "active-model": "4b-it-int4-fq-ov",
	//   "models": [
	//     {
	//       "name": "26b-q4-k-m-gguf",
	//       "description": "Test model description",
	//       "model-card-url": "https://example.com/model-card",
	//       "format": "GGUF",
	//       "quantization": "Q4_K_M",
	//       "capabilities": [
	//         "text"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-26b-a4b-q4-k-m-gguf",
	//         "mmproj-26b-bf16-gguf"
	//       ],
	//       "compatible-engines": [
	//         "arm-neon",
	//         "cpu",
	//         "cpu-avx1",
	//         "cpu-avx2",
	//         "cpu-avx512",
	//         "cuda-generic",
	//         "rocm-generic"
	//       ]
	//     },
	//     {
	//       "name": "30b-a3b-q4-k-m-gguf",
	//       "description": "Test model description",
	//       "model-card-url": "https://example.com/model-card",
	//       "format": "GGUF",
	//       "quantization": "Q4_K_M",
	//       "capabilities": [
	//         "text",
	//         "vision",
	//         "audio",
	//         "tool"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-30b-a3b-q4-k-m-gguf-1-of-6",
	//         "model-30b-a3b-q4-k-m-gguf-2-of-6",
	//         "model-30b-a3b-q4-k-m-gguf-3-of-6",
	//         "model-30b-a3b-q4-k-m-gguf-4-of-6",
	//         "model-30b-a3b-q4-k-m-gguf-5-of-6",
	//         "model-30b-a3b-q4-k-m-gguf-6-of-6"
	//       ],
	//       "compatible-engines": [
	//         "cpu",
	//         "cuda-generic",
	//         "rocm-generic"
	//       ]
	//     },
	//     {
	//       "name": "4b-it-int4-fq-ov",
	//       "alias": "4b-it",
	//       "description": "OpenVino 4b test model",
	//       "model-card-url": "https://example.com/model-card",
	//       "format": "OpenVINO_IR",
	//       "quantization": "int4-fq",
	//       "capabilities": [
	//         "text"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-4b-it-int4-fq-ov"
	//       ],
	//       "compatible-engines": [
	//         "intel-cpu",
	//         "intel-gpu",
	//         "intel-npu"
	//       ]
	//     }
	//   ]
	// }
}
