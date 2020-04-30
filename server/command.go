package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// annotationTagIsAllow ensure the tags received from Mattermost slash command are present
// into the tags list allowed
func annotationTagIsAllow(tag string, tags []string) bool {
	if len(tags) < 1 {
		return true
	}
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func createAnnotation(configuration *configuration, text string, tags []string) string {
	annotation := sdk.CreateAnnotationRequest{
		Text: text,
		Tags: tags,
	}
	ctx := context.Background()
	c := sdk.NewClient(configuration.GrafanaURL, configuration.GrafanaAPIKey, sdk.DefaultHTTPClient)
	result, _ := c.CreateAnnotation(ctx, annotation)
	return *result.Message
}

const commandHelp = `* |/fyi annotate "reason text" #tag1 #tag2 | - Create Grafana annotations
* |/fyi settings| - Display the configurations
`

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "fyi",
		DisplayName:      "For Your Information",
		Description:      "Create Grafana annotations.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: help, annotate, settings",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) getCommandResponse(args *model.CommandArgs, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         text,
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	config := p.getConfiguration()
	user, _ := p.API.GetUser(args.UserId)
	channel, _ := p.API.GetChannel(args.ChannelId)

	var (
		split      = strings.Fields(args.Command)
		command    = split[0]
		action     string
		parameters []string
	)
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}
	if command != "/fyi" {
		return &model.CommandResponse{}, nil
	}

	if action == "help" || action == "" {
		text := "###### Mattermost FYI Plugin - Slash Command Help\n" + strings.Replace(commandHelp, "|", "`", -1)
		return p.getCommandResponse(args, text), nil
	}

	switch action {
	case "annotate":
		if len(parameters) == 0 {
			text := fmt.Sprintf("Please specify a reason.\n\n------\n_%v_ \n\n-----\n%v", args.Command, strings.Replace(commandHelp, "|", "`", -1))
			return p.getCommandResponse(args, text), nil
		}
		tagsAnnotation := []string{"fyi"}
		textAnnotation := []string{}
		tagsConfiguration := strings.Split(config.Tags, ",")
		for _, field := range parameters {
			if len(field) > 1 && strings.Contains(field, "#") {
				if !annotationTagIsAllow(field[1:], tagsConfiguration) {
					text := fmt.Sprintf("Unknown tag **%v**, these tags are available \n _```%v```_ \n\n -----\n_%v_ \n\n -----\n%v", field, tagsConfiguration, args.Command, strings.Replace(commandHelp, "|", "`", -1))
					return p.getCommandResponse(args, text), nil
				}
				tagsAnnotation = append(tagsAnnotation, field[1:])
			} else {
				textAnnotation = append(textAnnotation, field)
			}
		}

		message := fmt.Sprintf("Reason: %v \nWho: %v\nFrom: %v", strings.Join(textAnnotation, " "), user.Username, channel.DisplayName)
		response := createAnnotation(p.configuration, message, tagsAnnotation)
		return p.getCommandResponse(args, response), nil
	case "settings":
		text := fmt.Sprintf("Grafana: %s\nTags: %s", config.GrafanaURL, config.Tags)
		return p.getCommandResponse(args, text), nil
	default:
		return p.getCommandResponse(args, "Unknown action, please use `/fyi help` to see all actions available."), nil
	}
}
