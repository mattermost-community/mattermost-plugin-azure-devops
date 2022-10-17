# Setting up your AzureDevops for oAuth

## Register your OAuth application on Azure DevOps
  - Log in to [AzureDevops](https://app.vsaex.visualstudio.com).
  - Click on the "Create new application" button in the bottom left corner.
  - Fill in details like Company name, Application name.
    ![Screenshot from 2022-07-29 13-24-47](https://user-images.githubusercontent.com/100013900/181712515-89efdb18-0f51-4194-b954-c0edb4188423.png)
  - For the "Application website" put the link of your Mattermost server, for example: If the Mattermost server is running on `https://<your-mattermost-url>` then the value would be `https://<your-mattermost-url>` and "Authorization callback URL" will be `https://<your-mattermost-url>/plugins/mattermost-plugin-azure-devops/api/v1/oauth/complete`

    ![Screenshot from 2022-07-29 13-25-38](https://user-images.githubusercontent.com/100013900/181712472-d4faec27-a61c-4565-9e27-fc7156241b17.png)

  - Under "Authorized scopes" check "Work items(full)" and "Code(full)".

    ![Screenshot from 2022-10-07 12-54-13](https://user-images.githubusercontent.com/55234496/194496403-dfd54566-ae6b-4daa-96c9-6a6c7e24c296.png)

  - Click on "Create application".
  - On successful creation, you will be navigated to a page having App ID, App Secret, and Client Secret

    ![Screenshot from 2022-10-07 13-02-54](https://user-images.githubusercontent.com/55234496/194498023-51eca666-7d58-47bf-80cf-74b4eecc1f04.png)

**Note**: This plugin uses the OAuth authentication protocol, and requires Third-party application access via OAuth for the organization to be enabled. To enable this setting, navigate to Organization Settings > Security > Policies, and set the Third-party application access via OAuth for the organization setting to On.

![image](https://user-images.githubusercontent.com/72438220/195812872-d97c6a80-2e84-4943-a1c4-3e570b10f995.png)
