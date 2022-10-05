export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';

export const eventTypeMap: Record<EventType, string> = {
    'workitem.created': 'Work Item Created',
    'workitem.updated': 'Work Item Updated',
    'workitem.deleted': 'Work Item Deleted',
    'workitem.commented': 'Work Item Commented On',
    'git.pullrequest.created': 'Pull Request Created',
    'git.pullrequest.updated': 'Pull Request Updated',
    'ms.vss-code.git-pullrequest-comment-event': 'Pull Request Commented On',
    'git.pullrequest.merged': 'Pull Request Merge Attempted',
    'git.push': 'Code Pushed',
};

export const serviceTypeMap: Record<EventType, string> = {
    'workitem.created': 'Boards',
    'workitem.updated': 'Boards',
    'workitem.deleted': 'Boards',
    'workitem.commented': 'Boards',
    'git.pullrequest.created': 'Repos',
    'git.pullrequest.updated': 'Repos',
    'ms.vss-code.git-pullrequest-comment-event': 'Repos',
    'git.pullrequest.merged': 'Repos',
    'git.push': 'Repos',
};

export const boards = 'boards';
export const repos = 'repos';

export const defaultPage = 0;
export const defaultPerPageLimit = 10;

export const SubscriptionFilterCreatedBy = {
    me: 'me',
    anyone: 'anyone',
};
