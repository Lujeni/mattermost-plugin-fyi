// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.github.manland.mattermost-plugin-fyi",
  "name": "FYI",
  "description": "Mattermost plugin to create Grafana annotation.",
  "homepage_url": "https://github.com/Lujeni/mattermost-plugin-fyi",
  "support_url": "https://github.com/Lujeni/mattermost-plugin-fyi",
  "icon_path": "assets/icon.svg",
  "version": "1.0.0",
  "min_server_version": "5.10.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "settings_schema": {
    "header": "To set up the Grafana Annotation plugin, you need to register a Grafana API KEY https://grafana.com/docs/grafana/latest/http_api/auth/.",
    "footer": "To report an issue, make a suggestion or a contribution, [check the repository](https://github.com/Lujeni/fyi).",
    "settings": [
      {
        "key": "GrafanaURL",
        "display_name": "Grafana URL",
        "type": "text",
        "help_text": "The base URL for using the plugin with a Grafana installation. Examples: https://grafana.example.com",
        "placeholder": "https://grafana.com",
        "default": null
      },
      {
        "key": "GrafanaAPIKey",
        "display_name": "Grafana API Key",
        "type": "text",
        "help_text": "The client ID for the OAuth app registered with GitLab.",
        "placeholder": "eyJrIjoiabUdzREMwcDEiLCJuIjoiZnlpIiwiaWQiOjF9",
        "default": null
      },
      {
        "key": "Tags",
        "display_name": "Grafana Annotations Tags",
        "type": "text",
        "help_text": "(Optional) Allow this Grafana annotations tags (strict mode).",
        "placeholder": "business,outage,reboot",
        "default": null
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
