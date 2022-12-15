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

This plugin integrates the services of Azure DevOps in Mattermost. For a stable production release, please download the latest version from the [Github Releases](https://github.com/mattermost/mattermost-plugin-azure-devops/releases) and follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

## Features

This plugin contains the following features:
- Right-hand sidebar (RHS) shows the list of linked projects and subscriptions created for each project with separate buttons to connect an account, link or unlink a project and add or delete a subscription.

- Preview of the work item, pull request, release or build URL: A preview of the work item, pull request, release or build for a linked project will be created when their respective URLs are posted in a channel and the user is connected to his/her Azure DevOps account.

- OAuth: A user can connect or disconnect to their Azure DevOps account using the slash command below or clicking on the "Connect Your Account" button in RHS.

    ```
    /azuredevops connect
    /azuredevops disconnect
    ```

- Link projects: A user can link a project existing on Azure DevOps using the slash command below or clicking on the "Link new project" button in RHS.

    ```
    /azuredevops link [project link]
    ```

- Unlink projects: A user can unlink a project appearing in the RHS under "Linked Projects" by clicking on the unlink-icon button.

- Create work items: A work item can be created using the slash command below.

    ```
    /azuredevops boards create [title] [description]
    ```
    On successful creation of a work item, you will get a message from the bot with the details of the newly created work item.

- Add subscriptions: A user can create subscriptions for a linked project to get notifications in a selected channel for selected events on work items, pull requests and pipelines.
To add a new subscription for a linked project click on the project title under "Linked Projects" in RHS then click on the "Add new subscription" button in the subscription view. Users can also create subscriptions using the slash command below.
    - For creating Boards subscriptions

    ```
    /azuredevops boards subscription add
    ```

    - For creating Repos subscriptions

    ```
    /azuredevops repos subscription add
    ```

    - For creating Pipelines subscriptions

    ```
    /azuredevops pipelines subscription add
    ```

    **Note:** Only Mattermost users who are project admins or team admins on the linked Azure DevOps project can create/delete a subscription.

- View/List subscriptions: A user can view the list of subscriptions for a project by going to the subscriptions list page after clicking on the project title under "Linked Projects" in the right-hand sidebar. Users can also view the list of all subscriptions for a channel by using the below slash command in the channel.

    - For listing Boards subscriptions

    ```
    /azuredevops boards subscription list anyone all_channels
    ```

    - For listing Repos subscriptions

    ```
    /azuredevops repos subscription list anyone all_channels 
    ```

    - For listing Pipelines subscriptions

    ```
    /azuredevops pipelines subscription list anyone all_channels 
    ```

    Supported filters on the above slash command:
    - CreatedBy: `me`(show all subscriptions created by the current Mattermost user), `anyone`(show all subscriptions created by any Mattermost user)
    - Show for all channels: When the filter `all_channels` is passed in the slash command then subscriptions for all channels are listed. You can skip this filter param to list the subscriptions of the current channel only.

    **Note:** Only Mattermost users who are project admins or team admins on the linked Azure DevOps project can view/list subscriptions that exist in a channel where they are not a member.

- Approve pipeline requests from your channel: 
  - A user can approve pipeline requests from within the channel by subscribing to "e Release deployment approval pending" or "Run stage waiting for approval" (for YAML pipelines) event notifications.

    ![image](https://user-images.githubusercontent.com/72438220/202639513-8fa58db0-bce9-46ed-864b-ebbae468a4b6.png)

  - When the running of a stage or a release deployment is awaiting approval, a notification card is posted in the channel with the option to approve or reject the request. Reviewing the request's specifics in the notification, approvers can take the necessary action.

    ![image](https://user-images.githubusercontent.com/72438220/202639612-1742c832-9a16-49a5-96ae-5cf6010f90d6.png)

  - Every check and approval scenario found in the Azure Pipelines interface is supported by the plugin, including single approver, multiple approvers (any one person, any order, in sequence), and teams as approvers.

- Delete subscriptions: A user can delete subscriptions for a project from RHS by going to the subscriptions list page after clicking on the project title under "Linked Projects". Users can also delete a subscription for a project by using the slash command below.

    - For deleting Boards subscriptions

    ```
    /azuredevops boards subscription delete [subscription id]
    ```

    - For deleting Repos subscriptions

    ```
    /azuredevops repos subscription delete [subscription id]
    ```

    - For deleting Pipelines subscriptions

    ```
    /azuredevops pipelines subscription delete [subscription id]
    ```

    **Note:** Only Mattermost users who are project admins or team admins on the linked Azure DevOps project can create/delete a subscription.

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/mattermost/mattermost-plugin-azure-devops/releases) and download the latest release for your Mattermost server.
2. Upload this file to the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. Enable the plugin from **System Console > Plugins > Mattermost Azure Devops plugin**.

## Setup

  - [Developer setup](./docs/developer_docs.md)
  - [OAuth app registration for Azure DevOps](./docs/oauth_setup.md)
  - [Plugin Setup](./docs/plugin_setup.md)

## Connecting to Azure Devops
  - Enter slash command `/azuredevops connect`.
  - You will get a response with a link to connect your Azure DevOps account.

    ![image](https://user-images.githubusercontent.com/100013900/181709568-9468b4a7-aaef-45a5-8968-882d560f43c3.png)

  - Click on that link. If it asks for a login, enter your credentials and connect to your account.

After connecting successfully, you will get a direct message from the Azure DevOps bot containing a Welcome message and some useful information. 

**Note:** You will only get a direct message from the bot if your Mattermost server is configured to allow direct messages between any users on the server. If your server is configured to allow direct messages only between two users of the same team, then you will not get any direct messages.
