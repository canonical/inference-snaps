package snap_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ComponentSizes() (map[string]int64, error) {
	components, err := componentsOfCurrentSnap()
	if err != nil {
		return nil, fmt.Errorf("error finding components of current snap: %v", err)
	}

	componentSizes := make(map[string]int64)
	for _, component := range components {
		componentSizes[component.Name] = component.Download.Size
	}

	return componentSizes, nil
}

func componentsOfCurrentSnap() ([]snapResources, error) {
	snapName := os.Getenv("SNAP_NAME")
	if snapName == "" {
		return nil, fmt.Errorf("SNAP_NAME is not set. Likely not inside a snap")
	}

	snapRevisionStr := os.Getenv("SNAP_REVISION")
	if snapRevisionStr == "" {
		return nil, fmt.Errorf("SNAP_REVISION is not set")
	}

	if strings.HasPrefix(snapRevisionStr, "x") {
		return nil, fmt.Errorf("not installed from store")
	}

	snapRevision, err := strconv.ParseInt(snapRevisionStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing snap revision: %v", err)
	}

	info, err := snapInfo(snapName)
	if err != nil {
		return nil, fmt.Errorf("error getting snap info: %v", err)
	}

	components, err := snapComponents(info.SnapId, int(snapRevision), os.Getenv("SNAP_ARCH"))
	if err != nil {
		return nil, fmt.Errorf("error getting components: %v", err)
	}

	return components, nil
}

// snapInfo fetches information about a respective snap from the store, based on its name
// Docs: https://api.snapcraft.io/docs/info.html
func snapInfo(snapName string) (*snapInfoResponse, error) {
	url := "https://api.snapcraft.io/v2/snaps/info/" + snapName

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new http request: %v", err)
	}

	req.Header.Add("Snap-Device-Series", "16")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status not OK: %d", resp.StatusCode)
	}

	info := snapInfoResponse{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &info, nil
}

func snapComponents(snapId string, revision int, snapArch string) ([]snapResources, error) {
	/*
		From sniffing the traffic between snapd and the snap store, we see the refresh endpoint is used to look up
		available components for the respective snap, their revisions, their sizes, and their download URLs.
	*/
	refreshResponse, err := snapRefresh(snapId, revision, snapArch)
	if err != nil {
		return nil, fmt.Errorf("error fetching refresh data from store: %v", err)
	}

	if len(refreshResponse.Results) == 0 {
		return nil, fmt.Errorf("store returned no refresh results")
	}

	for _, result := range refreshResponse.Results {
		if result.SnapId == snapId {
			return result.Snap.Resources, nil
		}
	}

	return nil, fmt.Errorf("no refresh results found for snap id %s", snapId)
}

// snapRefresh makes a call to the store to fetch refresh information
// Docs: https://api.snapcraft.io/docs/refresh.html
func snapRefresh(snapId string, revision int, snapArch string) (*snapRefreshResponse, error) {
	request := snapRefreshRequest{
		Context: []snapRefreshContext{
			{
				SnapId:          snapId,
				InstanceKey:     snapId,
				Revision:        revision,
				TrackingChannel: "",
			}},
		Actions: []snapRefreshActions{
			{
				Action:      "refresh",
				InstanceKey: snapId,
				SnapId:      snapId,
				Revision:    revision,
			},
		},
		Fields: []string{"resources"},
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request to json: %v", err)
	}

	url := "https://api.snapcraft.io/v2/snaps/refresh"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, fmt.Errorf("error creating new http request: %v", err)
	}

	req.Header.Add("Snap-Device-Series", "16")
	req.Header.Add("Snap-Device-Architecture", snapArch)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status not OK: %d", resp.StatusCode)
	}

	var response snapRefreshResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &response, nil
}
