export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';

export enum ToggleSwitchLabelPositioning {
    Left = 'left',
    Right = 'right',
}

export const ToggleLabel = 'Show all subscriptions';

export const boardsEventTypeMap: Record<EventType, string> = {
    create: 'Work Item Created',
    update: 'Work Item Updated',
    delete: 'Work Item Deleted',
};

export const channelType = {
    priivate: 'P',
};
