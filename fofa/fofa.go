package fofa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://fofa.info/api/v1/"

type Client struct {
	Email string
	Key   string
}

func NewClient(email, key string) *Client {
	return &Client{
		Email: email,
		Key:   key,
	}
}

func (c *Client) GetData(uri string, params string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s?email=%s&key=%s&%s", apiURL, uri, c.Email, c.Key, params)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	if isError, ok := result["error"].(bool); ok && isError {
		return nil, fmt.Errorf("%s", result["errmsg"])
	}

	return result, nil
}

func (c *Client) GetAccountInfo() (map[string]interface{}, error) {
	return c.GetData("info/my", "")
}

func (c *Client) SearchData(rule string, page int, size int) (map[string]interface{}, error) {
	ruleBase64 := base64.URLEncoding.EncodeToString([]byte(rule))
	params := fmt.Sprintf("qbase64=%s&page=%d&size=%d&fields=ip,port,country,as_organization", ruleBase64, page, size)
	return c.GetData("search/all", params)
}
