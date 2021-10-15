### mattermost-plugin-outline
Mattermost **[Outline](https://www.getoutline.com/)** plugin allows you to search your teams documents.

![mattermost](https://raw.githubusercontent.com/Lujeni/mattermost-plugin-outline/main/assets/mattermost.png)

## Installation
In Mattermost 5.16 and later, this plugin is included in the Plugin Marketplace which can be accessed from **Main Menu > Plugins Marketplace**. You can install the plugin and then configure it via the [Plugin Marketplace "Configure" button](#configuration).

In Mattermost 5.13 and earlier, follow these steps:
1. Go to https://github.com/Lujeni/mattermost-plugin-outline/releases to download the latest release file in zip or tar.gz format.
2. Upload the file through **System Console > Plugins > Management**, or manually upload it to the Mattermost server under plugin directory. See [documentation](https://docs.mattermost.com/administration/plugins.html#set-up-guide) for more details.

## Configuration
### Step 1: Retrieve outline API Key
   
1. Go to https://www.getoutline.com/developers#section/Authentication

### Step 2: Configure plugin in Mattermost

1. Go to **System Console > Plugins > Outline** and fill the form

## Development
This plugin contains only a server portion (no web app).

Use make to build distributions of the plugin that you can upload to a Mattermost server.

```
$ make
