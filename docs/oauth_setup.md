# Setting up your AzureDevops for oAuth

## Register your OAuth application on Azure DevOps
  - Log in to [AzureDevops](https://app.vsaex.visualstudio.com).
  - Click on the "Create new application" button in the bottom left corner.
  - Fill in details like Company name, Application name.
    ![Screenshot from 2022-07-29 13-24-47](https://user-images.githubusercontent.com/100013900/181712515-89efdb18-0f51-4194-b954-c0edb4188423.png)
  - For the "Application website" put the link of your Mattermost server, for example: If the Mattermost server is running on `https://<your-mattermost-url>` then the value would be `https://<your-mattermost-url>` and "Authorization callback URL" will be `https://<your-mattermost-url>/plugins/mattermost-plugin-azure-devops/api/v1/oauth/complete`

    ![Screenshot from 2022-07-29 13-25-38](https://user-images.githubusercontent.com/100013900/181712472-d4faec27-a61c-4565-9e27-fc7156241b17.png)

  - Under "Authorized scopes" check "Work items(full)" and "Code(read)".

    ![Screenshot from 2022-07-29 13-26-00](https://user-images.githubusercontent.com/100013900/181712419-3251837c-ea66-47b9-8b5b-7ae950816f2d.png)

  - Click on "Create application".
  - On successful creation, you will be navigated to a page having App ID, App Secret, and Client Secret

    ![Screenshot from 2022-07-29 13-26-16](https://user-images.githubusercontent.com/100013900/181712321-e049ce0c-4b22-4c35-a60e-123d24fb0791.png)
