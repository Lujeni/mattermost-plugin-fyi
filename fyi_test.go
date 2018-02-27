package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	var config Config
	err := envconfig.Process("fyi", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleIndex(w, r, config)
	}).Methods("GET", "POST")

	return r
}

func SimulateRequest(request *http.Request) (*httptest.ResponseRecorder, CommandResponse) {
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	commandResponse := CommandResponse{}
	if err := json.NewDecoder(response.Body).Decode(&commandResponse); err != nil {
		log.Fatalln(err)
	}
	return response, commandResponse
}

// TestNotAllowedMethod send unsupported request (GET)
func TestNotAllowedMethod(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Equal(t, "GET method not allowed, use POST", responseJSON.Text)
}

// TestEmptyForm send a request without body
func TestEmptyForm(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Equal(t, "Unable to parse form :: missing form body", responseJSON.Text)
}

// TestEmptyText send a request with an empty text form value
func TestEmptyText(t *testing.T) {
	data := url.Values{}
	data.Set("text", "")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "No message specify")
}

// TestPayloadText send a text without the grafana configuration
func TestPayloadText(t *testing.T) {
	data := url.Values{}
	data.Set("text", "foo")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "Unable to post grafana")
}

// TestPayloadWithoutTags send a text without a tag
func TestPayloadWithoutTags(t *testing.T) {
	os.Setenv("FYI_TAGS", "infra,outage,marketing")
	data := url.Values{}
	data.Set("text", "foo")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "No tag specify")
}

// TestPayloadWithBadTags send a text with a bad tag
func TestPayloadWithBadTags(t *testing.T) {
	os.Setenv("FYI_TAGS", "infra,outage,marketing")
	data := url.Values{}
	data.Set("text", "foo #foo")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "Unknown tag")
}

// TestPayloadWithTags send a text with a good tagw
func TestPayloadWithTags(t *testing.T) {
	os.Setenv("FYI_TAGS", "infra,outage,marketing")
	data := url.Values{}
	data.Set("text", "foo #infra")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "Unable to post grafana")
}

// TestPayloadWithBadToken send a text with a (bad) token verification
func TestPayloadWithBadToken(t *testing.T) {
	os.Setenv("FYI_TOKEN", "foobar")
	data := url.Values{}
	data.Set("text", "foo #infra")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "Bad token received")
}

// TestPayloadWithToken send a text with a token verification
func TestPayloadWithToken(t *testing.T) {
	os.Setenv("FYI_TOKEN", "foobar")
	data := url.Values{}
	data.Set("text", "foo #infra")
	data.Set("token", "foobar")
	request, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	responseHTTP, responseJSON := SimulateRequest(request)

	assert.Equal(t, 200, responseHTTP.Code)
	assert.Contains(t, responseJSON.Text, "Unable to post grafana")
}
