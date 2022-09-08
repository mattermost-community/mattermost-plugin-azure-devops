# Mattermost Azure DevOps Plugin
## Table of Contents
- [License](#license)
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Setup](#setup)
- [Connecting to Azure DevOps](#connecting-to-azure-devops)
- [References](#references)
- [Development](#development)

## License

See the [LICENSE](./LICENSE) file for license rights and limitations.

## Overview

This plugin integrates the services of Azure DevOps in Mattermost. For a stable production release, please download the latest version from the [Github Releases](https://github.com/Brightscout/mattermost-plugin-azure-devops/releases) and follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

## Features

This plugin contains the following features:
- Right-hand sidebar (RHS) shows the list of linked projects and subscriptions created for each project with separate buttons to connect an account, link or unlink a project and add or delete a subscription.

- Preview work item: A preview of the work item will be created when a work item URL is posted in a channel and the user account is connected to his Azure DevOps account.

- OAuth: A user can connect or disconnect to their Azure DevOps account using the slash command below or clicking on the "Connect Your Account" button in RHS.

    ```
    - /azuredevops connect
    - /azuredevops disconnect
    ```

- Link projects: A user can link a project existing on Azure DevOps using the slash command below or clicking on the "Link new project" button in RHS.

    ```
    - /azuredevops link [project link]
    ```

- Unlink projects: A user can unlink a project appearing in the RHS under "Linked Projects" by clicking on the unlink-icon button.

- Create work items: A work item can be created using the slash command below.

    ```
    - /azuredevops boards create
    ```
    On successful creation of a work item, you will get a message from the bot with the details of the newly created work item.

- Add subscriptions: A user can create subscriptions for a linked project to get messages in a selected channel for selected events on work items like Create, Update and Delete.
To add a new subscription for a linked project click on the project title under "Linked Projects" in RHS then click on the "Add new subscription" button in the subscription view. Users can also create subscriptions using the slash command below.

    ```
    - /azuredevops boards subscribe
    ```

- View/List subscriptions: A user can view the list of subscriptions for a project by going on the subscriptions list page after clicking on the project title under "Linked Projects". Users can also view the list of all subscriptions for a channel by using the below slash command in the channel.

    ```
    - /azuredevops boards subscriptions
    ```

- Delete subscriptions: A user can delete subscriptions for a project from RHS by going on the subscriptions list page after clicking on the project title under "Linked Projects". Users can also delete a subscription for a project by using the slash command below.

    ```
    - /azuredevops boards unsubscribe [subscription id]
    ```

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-azure-devops/releases) and download the latest release for your Mattermost server.
2. Upload this file to the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. Enable the plugin from **System Console > Plugins > Mattermost Azure Devops plugin**.

## Setup

  - [Developer setup](./developer_docs.md)
  - [OAuth app registration for Azure DevOps](./oauth_setup.md)
  - [Plugin Setup](./plugin_setup.md)

## Connecting to Azure Devops
  - Enter slash command `/azuredevops connect`.
  - You will get a response with a link to connect your Azure DevOps account.
  ![Screenshot from 2022-07-29 13-01-14](https://user-images.githubusercontent.com/100013900/181709568-9468b4a7-aaef-45a5-8968-882d560f43c3.png)
  - Click on that link. If it asks for a login, enter your credentials and connect to your account.

After connecting successfully, you will get a direct message from the Azure DevOps bot containing a Welcome message and some useful information. 

**Note:** You will only get a direct message from the bot if your Mattermost server is configured to allow direct messages between any users on the server. If your server is configured to allow direct messages only between two users of the same team, then you will not get any direct messages.

## References
You can read below mentioned documents to get knowledge about the Azure DevOps Rest API services.

- [OAuth](https://docs.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/oauth)
- [Work Items](https://docs.microsoft.com/en-us/rest/api/azure/devops/wit/work-items)
- [Subscriptions](https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/subscriptions)

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
