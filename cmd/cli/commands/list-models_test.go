package commands

import (
	"fmt"
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func prepareModelsTestData() (*listModelsCommand, *outputModels, error) {
	cache := storage.NewMockCache()
	err := cache.SetActiveModel("4b-it-int4-fq-ov")
	if err != nil {
		return nil, nil, fmt.Errorf("error setting active model name: %v", err)
	}

	manifests, err := models.LoadManifests("../../../test_data/models")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading models: %v", err)
	}

	var allModels []common.ModelDetails
	for _, manifest := range manifests {
		details, err := common.NewModelDetails(&manifest)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating model details for %s: %v", manifest.ID, err)
		}
		allModels = append(allModels, details)
	}

	ctx := &common.Context{
		ModelsDir: "../../../test_data/models",
		Cache:     cache,
		Config:    nil,
	}
	cmd := listModelsCommand{Context: ctx}

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

func TestListModelsJson(t *testing.T) {
	cmd, modelsList, err := prepareModelsTestData()
	if err != nil {
		t.Fatalf("Error preparing test data: %v", err)
	}

	err = cmd.printModelsJson(*modelsList)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListModelsTable(t *testing.T) {
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

	expectedTable := `NAME     CAPABILITIES               DISK SIZE                                   
26b      text                       6G                                          
30b-a3b  text, vision, audio, tool  6G                                          
4b*      text                       6G                                          
`

	if tableStr != expectedTable {
		t.Errorf("Models table not as expected.\n\nGot:\n\n%s\n\nWant:\n\n%s", tableStr, expectedTable)
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
		if m.ID == "4b-it-int4-fq-ov" {
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
	//       "id": "4b-it-int4-fq-ov",
	//       "name": "4b",
	//       "description": "OpenVino 4b test model",
	//       "model-card-url": "https://example.com/model-card",
	//       "quantization": "int4-fq",
	//       "capabilities": [
	//         "text"
	//       ],
	//       "disk-size": "6G",
	//       "components": [
	//         "model-4b-it-int4-fq-ov"
	//       ]
	//     }
	//   ]
	// }
}
