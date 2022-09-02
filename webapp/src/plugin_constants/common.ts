// Plugin configs
export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';

export const boardsEventTypeMap: Record<eventType, string> = {
    create: 'Work Item Created',
    update: 'Work Item Updated',
    delete: 'Work Item Deleted',
};
