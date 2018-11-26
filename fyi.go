package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/kelseyhightower/envconfig"
)

// Config represents all variables to run the command
type Config struct {
	Debug         bool     `default:"true"`
	Host          string   `default:"0.0.0.0"`
	Port          int      `default:"8888"`
	Token         string   `required:"false"`
	Tags          []string `required:"false"`
	GrafanaHost   string   `envconfig:"grafana_host"`
	GrafanaAPIKey string   `envconfig:"grafana_api_key"`
	Username      string   `default:"ForYourInformation"`
	IconURL       string   `default:"https://avatars2.githubusercontent.com/u/757902?s=460&v=4"`
}

type OutgoingWebhookPayload struct {
	Token       string `json:"token"`
	TeamId      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Timestamp   int64  `json:"timestamp"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	PostId      string `json:"post_id"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
	FileIds     string `json:"file_ids"`
}

type CommandResponse struct {
	IconURL      string `json:"icon_url"`
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
	Username     string `json:"username"`
}

type GrafanaAnnotation struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Time string   `json:"time"`
}

type GrafanaAnnotationResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
	EndID   int    `json:"endId"`
}

// TagIsAllow ensure the tags received from Mattermost slash command are present
// into the tags list allowed
func TagIsAllow(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// SendGrafanaAnnotation prepare the annotation and send it to Grafana API
func SendGrafanaAnnotation(config Config, annotation GrafanaAnnotation) (string, error) {
	response := GrafanaAnnotationResponse{}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(annotation)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/api/annotations", config.GrafanaHost), b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", config.GrafanaAPIKey))
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Unable to post grafana annotation %v", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error in grafana server response %v", res.Status)
	}
	json.NewDecoder(res.Body).Decode(&response)
	return response.Message, nil
}

func ExtractDate(text string) []string {
	r := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}){1,2}`)
	return r.FindAllString(text, 2)
}

// ProcessCommand do all the stuff, sanity check, send annotations ect...
func ProcessCommand(r *http.Request, response *CommandResponse, config Config) string {
	if err := r.ParseForm(); err != nil {
		msg := fmt.Sprintf("Unable to parse form :: %v", err.Error())
		log.Printf(msg)
		return msg
	}

	payload := new(OutgoingWebhookPayload)
	decoder := schema.NewDecoder()
	err := decoder.Decode(payload, r.PostForm)

	if err != nil {
		msg := fmt.Sprintf("Unable to decode struct :: %v", err.Error())
		log.Printf(msg)
		return msg
	}

	if config.Token != "" && config.Token != payload.Token {
		msg := fmt.Sprintf("Bad token received :: %v", payload.Token)
		log.Printf(msg)
		return msg
	}
	text := strings.Fields(payload.Text)
	dates := ExtractDate(payload.Text)
	if len(dates) > 0 {
		layout := "2006-01-01 01:00:00"
		t, err := time.Parse(layout, dates[0])
		if err != nil {
			msg := fmt.Sprintf("Unable to parse time :: %v", dates[0])
			log.Printf(msg)
			return msg
		}
		fmt.Printf("OK %v", t)
	}
	tagsAnnotation := []string{"fyi", payload.UserName}
	textAnnotation := []string{}

	for _, field := range text {
		if len(field) > 1 && strings.Contains(field, "#") {
			if !TagIsAllow(field[1:], config.Tags) {
				return fmt.Sprintf("Unknown tag **%v**, these tags are available \n - tags: _```%v```_", field, config.Tags)
			}
			tagsAnnotation = append(tagsAnnotation, field[1:])
		} else {
			textAnnotation = append(textAnnotation, field)
		}
	}

	if len(tagsAnnotation) <= 2 && len(config.Tags) > 0 {
		return fmt.Sprintf("No tag specify, **one** of these tags are mandatory \n - tags: _```%v```_", config.Tags)
	}

	if len(textAnnotation) == 0 {
		return fmt.Sprintf("No message specify \n - example: _```/command reboot server #outage```_")
	}

	grafanaAnnotation := GrafanaAnnotation{Text: strings.Join(textAnnotation[:], " "), Tags: tagsAnnotation}
	msg, err := SendGrafanaAnnotation(config, grafanaAnnotation)
	if err != nil {
		return err.Error()
	}
	return msg
}

// HandleHealthz is usefull for Liveness and Readiness Probes
func HandleHealthz(w http.ResponseWriter, r *http.Request, config Config) {
	// check grafana connectivity
	_, errGrafana := net.Dial("tcp", fmt.Sprintf("%v:80", config.GrafanaHost))
	if errGrafana != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - %v", errGrafana)))
	}

	// check local endpoint
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", config.Host, config.Port))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - local API - %v", err)))
	}
	return
}

func HandleIndex(w http.ResponseWriter, r *http.Request, config Config) {
	var response CommandResponse
	response.ResponseType = "ephemeral"
	response.Username = config.Username
	response.IconURL = config.IconURL
	if r.Method == "GET" {
		response.Text = "GET method not allowed, use POST"
	} else {
		response.Text = ProcessCommand(r, &response, config)
	}
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return
}

func main() {
	var config Config
	err := envconfig.Process("fyi", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	address := fmt.Sprintf("%v:%v", config.Host, config.Port)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleIndex(w, r, config)
	}).Methods("GET", "POST")
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		HandleHealthz(w, r, config)
	}).Methods("GET")

	log.Printf("ListenAndServe")
	log.Printf("|__Address ::	%v\n", address)
	log.Printf("|__Grafana :: %v\n", config.GrafanaHost)
	log.Printf("|__Tags :: %v\n", config.Tags)
	loggerHandler := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(address, loggerHandler))
}
