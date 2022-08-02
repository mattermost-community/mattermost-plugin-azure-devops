# Mattermost Azure DevOps Plugin
## Table of Contents
- [License](#license)
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Setup](#setup)
- [Connecting to Azure DevOps](#connecting-to-azure-devops)
- [Development](#development)

## License

See the [LICENSE](./LICENSE) file for license rights and limitations.

## Overview

This plugin integrates the services of Azure DevOps in Mattermost. For a stable production release, please download the latest version from the [Github Release](https://github.com/Brightscout/mattermost-plugin-azure-devops/releases) and follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

## Features

- oAuth: A user may connect or disconnect to their Azure DevOps account using the slash command below.

    ```
    - /azuredevops connect
    - /azuredevops disconnect
    ```

- Link a project: You can link a project existing on Azure DevOps using the slash command below.

    ```
    - /azuredevops link <project link>
    ```

- Unink a project: You can unlink a project appearing in the RHS under "Linked Projects".


- Right-hand sidebar (RHS) shows the list of projects linked to your current channel. Each project will have an option to **unlink** a project from the current channel.

- Preview work item: A preview of the work item will be created when a work item URL is posted in a channel.

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-azure-devops/releases) and download the latest release for your Mattermost server.
2. Upload this file to the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. Enable the plugin from **System Console > Plugins > Mattermost Azure Devops plugin**.

## Setup

  - [oAuth app registration for Azure DevOps](./docs/oauth_setup.md)
  - [Plugin Setup](./docs/plugin_setup.md)

## Connecting to Azure Devops
  - Enter slash command `/azuredevops connect`.
  - You will get a response with a link to connect your Azure DevOps account.
  ![Screenshot from 2022-07-29 13-01-14](https://user-images.githubusercontent.com/100013900/181709568-9468b4a7-aaef-45a5-8968-882d560f43c3.png)
  - Click on that link. If it asks for a login, enter your credentials and connect to your account.

## Development

### Setup

Make sure you have the following components installed:  

- Go - v1.16 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer to [Installing to a custom location](https://golang.org/doc/install#install).
- Make

### Building the plugin

Run the following command in the plugin repo to prepare a compiled, distributable plugin zip:

```bash
make dist
```

After a successful build, a `.tar.gz` file in `/dist` folder will be created which can be uploaded to Mattermost. To avoid having to manually install your plugin, deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    },
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

and then deploy your plugin:

```bash
make deploy
```

You may also customize the Unix socket path:

```bash
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a web app, watch for changes and deploy those automatically:

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
