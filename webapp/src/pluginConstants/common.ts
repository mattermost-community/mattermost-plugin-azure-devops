export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';
export const StatusCodeForbidden = 403;

export const deleteAllSubscriptionsMessage = 'Delete all your subscriptions associated with this project';
export const projectLinkedSuccessfullyMessage = 'Project linked successfully.';
export const projectAlreadyLinkedMessage = 'Project already linked.';

export const eventTypeBoards = {
    'workitem.created': 'Work Item Created',
    'workitem.updated': 'Work Item Updated',
    'workitem.deleted': 'Work Item Deleted',
    'workitem.commented': 'Work Item Commented On',
};

export const eventTypeReposKeys = {
    created: 'git.pullrequest.created',
    updated: 'git.pullrequest.updated',
    commented: 'ms.vss-code.git-pullrequest-comment-event',
    merged: 'git.pullrequest.merged',
    codePushed: 'git.push',
};

export const eventTypeRepos = {
    'git.pullrequest.created': 'Pull Request Created',
    'git.pullrequest.updated': 'Pull Request Updated',
    'ms.vss-code.git-pullrequest-comment-event': 'Pull Request Commented On',
    'git.pullrequest.merged': 'Pull Request Merge Attempted',
    'git.push': 'Code Pushed',
};

export const eventTypeMap: Record<EventType, string> = {
    ...eventTypeBoards,
    ...eventTypeRepos,
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
export const serviceType = 'serviceType';
export const eventType = 'eventType';

export const defaultPage = 0;
export const defaultPerPageLimit = 10;

export const subscriptionFilters = {
    createdBy: {
        me: 'me',
        anyone: 'anyone',
    },
    serviceType: {
        boards: 'boards',
        repos: 'repos',
        all: 'all',
    },
    eventType: {
        boards: {
            ...eventTypeBoards,
        },
        repos: {
            ...eventTypeRepos,
        },
        all: 'all',
    },
};

export const defaultSubscriptionFilters = {
    createdBy: subscriptionFilters.createdBy.anyone,
    serviceType: subscriptionFilters.serviceType.all,
    eventType: subscriptionFilters.eventType.all,
};

export const filterLabelValuePairAll = {
    value: 'all',
    label: 'All',
};
