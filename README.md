### mattermost-plugin-fyi
Mattermost plugin to easily create Grafana annotations.

![mattermost](https://raw.githubusercontent.com/Lujeni/fyi/master/assets/mattermost.png)

## Installation
In Mattermost 5.16 and later, the FYI plugin is included in the Plugin Marketplace which can be accessed from **Main Menu > Plugins Marketplace**. You can install the plugin and then configure it via the [Plugin Marketplace "Configure" button](#configuration).

In Mattermost 5.13 and earlier, follow these steps:
1. Go to https://github.com/Lujeni/mattermost-plugin-fyi/releases to download the latest release file in zip or tar.gz format.
2. Upload the file through **System Console > Plugins > Management**, or manually upload it to the Mattermost server under plugin directory. See [documentation](https://docs.mattermost.com/administration/plugins.html#set-up-guide) for more details.

## Configuration
### Step 1: Generate Grafana API Key
   
1. Go to https://grafana.com/docs/grafana/latest/http_api/auth/
2. Set the following values:
   - **Name**: FYI (or whatever)
   - **Role**: `Editor`
3. Copy the Token

### Step 2: Configure plugin in Mattermost

1. Go to **System Console > Plugins > FYI** and fill the form
3. Go to **Plugins Marketplace > FYI > Configure > Enable Plugin** and click **Enable** to enable the FYI plugin.

## Developping
This plugin contains only a server portion (no web app).

Use make to build distributions of the plugin that you can upload to a Mattermost server.

```
$ make
```
