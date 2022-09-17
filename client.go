package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Credentials struct {
	accessToken string
}

func (c *Credentials) GetAuthHeaderValue() string {
	return "Bearer " + c.accessToken
}

type Client struct {
	Credentials Credentials
}

func (c *Client) BuildClient(token string) *Client {
	c.Credentials = Credentials{
		accessToken: token,
	}
	return c
}

type MiroItemsResponse struct {
	Size     int16           `json:"size"`
	Limit    int16           `json:"limit"`
	Total    int16           `json:"total"`
	DataList []*MiroResponse `json:"data"`
	Links    *Links          `json:"links"`
	Type     string          `json:"type"`
}

func (items *MiroItemsResponse) ToJson(body []byte) error {
	err := json.Unmarshal(body, &items)
	if err != nil {
		return err
	}

	return nil
}

type MiroResponse struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Links      *Links    `json:"links"`
	CreatedAt  string    `json:"createdAt"`
	CreatedBy  *User     `json:"createdBy"`
	Data       *Data     `json:"data"`
	Geometry   *Geometry `json:"geometry"`
	ModifedAt  string    `json:"modifiedAt"`
	Modifiedby *User     `json:"modifiedBy"`
	Position   *Position `json:"position"`
}

type Links struct {
	Self string `json:"self"`
}

type User struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Data struct {
	ImageUrl string `json:"imageUrl"`
	Title    string `json:"title"`
}

type Geometry struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type Position struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Origin string  `json:"origin"`
}

/*
DoRequest performs a REST call
*/
func (c *Client) DoRequest(url string, body []byte) []byte {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", c.Credentials.GetAuthHeaderValue())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return output
}
