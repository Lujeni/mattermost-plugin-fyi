### fyi
Mattermost slash command to create Grafana annotations.

![example](https://raw.githubusercontent.com/Lujeni/fyi/master/example.png)

### Usage
#### Docker
```bash
# Create your Mattermost slash command
# Follow this steps: https://docs.mattermost.com/developer/slash-commands.html

# Pull
$ docker build . -t lujeni/fyi

# Run !
$ docker run -e FYI_GRAFANA_HOST=https://grafana.com -e FYI_GRAFANA_API_KEY=eyJrIjoibXppaTB5NXVu --rm -it lujeni/fyi
```

#### Kubernetes
```
$ kubectl apply -f kubernetes
```

#### Manual
```bash
# Create your Mattermost slash command
# Follow this steps: https://docs.mattermost.com/developer/slash-commands.html

# Install the project's dependencies
$ make deps

# Mattermost security token [optional]
$ export FYI_TOKEN=euhi1mfrhpbuny17qaq

# Tags allowed [optional]
$ export FYI_TAGS=infra,outage,marketing

# Setup your grafana API
$ export FYI_GRAFANA_HOST=https://grafana.com
$ export FYI_GRAFANA_API_KEY=eyJrIjoibXppaTB5NXVu

# Run !
$ make run
```

#### Pre-packaged
You don't have to install any other software.
Packages are available on the [releases page](http://github.com/Lujeni/fyi/releases).

#### Configuration
```go
type Config struct {
	Debug         bool     `default:"true"`
	Host          string   `default:"0.0.0.0"`
	Port          int      `default:"8888"`
	Token         string   `required:"false"`
	Tags          []string `required:"false"`
	GrafanaHost   string   `envconfig:"grafana_host"`
	GrafanaApiKey string   `envconfig:"grafana_api_key"`
	Username      string   `default:"ForYourInformation"`
	IconURL       string   `default:"https://avatars2.githubusercontent.com/u/757902?s=460&v=4"`
}
```

#### Tests
```bash
$ make test
```
