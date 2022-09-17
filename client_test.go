package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestDoRequest(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	boardId := "uXjVPXPlJik%3D"
	url := fmt.Sprintf("https://api.miro.com/v2/boards/%s/items?limit=10", boardId)

	// mock the request
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Accept") != "application/json" {
				t.Errorf("Expected Accept: application/json header, got: %s", req.Header.Get("Accept"))
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"size": 10, "limit": 10, "total": 172,
				"data": []*MiroResponse{
					{
						Id:   "123",
						Type: "image",
						Data: &Data{
							Title: "a thing",
						},
					},
				},
				"links": map[string]interface{}{
					"self": "https://api.miro.com/v2/boards/uXjVPXPlJik=/items?limit=10&cursor=",
					"next": "https://api.miro.com/v2/boards/uXjVPXPlJik=/items?limit=10&cursor=MzQ1ODc2NDUzMzQ3MTExODg0M34%3D",
					"last": "https://api.miro.com/v2/boards/uXjVPXPlJik=/items?limit=10&cursor=MzQ1ODc2NDUzMzQ4MzM0Mjk4M34%3D",
				},
				"type": "cursor-list",
			})
			return resp, err
		},
	)

	var client Client
	client.BuildClient("token")

	resp := client.DoRequest(url, nil)

	//info := httpmock.GetCallCountInfo()
	//assert.Equal(info[], 1, "The two words should be the same.")

	var items MiroItemsResponse
	err := items.ToJson(resp)
	if err != nil {
		fmt.Println(err)
	}

	// Assertions
	if items.DataList[0].Type != "image" {
		t.Fatalf(`Do Request: Name`)
	}
	if items.DataList[0].Id != "123" {
		t.Fatalf(`Do Request: Id`)
	}
}
