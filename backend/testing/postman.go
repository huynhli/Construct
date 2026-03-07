package testing

import (
	"backend/handlers"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const apiURL = "https://construct-pbbt.onrender.com/api/v1/"

var jwtToken, _ = handlers.GenerateJWT("TestUser", 999, true)

func TestGetWithJWT(t *testing.T) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode, "Expected status 200")
	assert.Contains(t, string(body), "expected_field_or_value")
}

func TestPostWithJWT(t *testing.T) {
	payload := map[string]interface{}{
		"title": "Test Task",
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, 201, resp.StatusCode, "Expected status 201")
	assert.Contains(t, string(body), "Test Task")
}
