# Setting up your AzureDevops for oAuth

## Register your OAuth application on Azure DevOps
  - Log in to [AzureDevops](https://app.vsaex.visualstudio.com).
  - Click on the "Create new application" button in the bottom left corner.
  - Fill in details like Company name, Application name.

    ![image](https://user-images.githubusercontent.com/100013900/181712515-89efdb18-0f51-4194-b954-c0edb4188423.png)

  - For the "Application website" put the link of your Mattermost server, for example: If the Mattermost server is running on `https://<your-mattermost-url>` then the value would be `https://<your-mattermost-url>` and "Authorization callback URL" will be `https://<your-mattermost-url>/plugins/mattermost-plugin-azure-devops/api/v1/oauth/complete`

    ![image](https://user-images.githubusercontent.com/55234496/204722294-64fce47d-8669-4bb7-8a69-d1e025129999.png)

  - Under "Authorized scopes" check "Work items(full)", Build(read and execute), Release(read, write, execute and manage) and "Code(full)".

    ![image](https://user-images.githubusercontent.com/55234496/204722689-b3feef42-0aa5-42ff-b011-efe1a6dad4c6.png)

  - Click on "Create application".
  - On successful creation, you will be navigated to a page having App ID, App Secret, and Client Secret

    ![image](https://user-images.githubusercontent.com/55234496/204722316-dddac330-616a-4a2a-a144-eb4e2611265b.png)

**Note**: This plugin uses the OAuth authentication protocol and requires Third-party application access via OAuth for the organization to be enabled. To enable this setting, navigate to Organization Settings > Security > Policies, and set the Third-party application access via OAuth for the organization setting to On.

![image](https://user-images.githubusercontent.com/72438220/195812872-d97c6a80-2e84-4943-a1c4-3e570b10f995.png)
