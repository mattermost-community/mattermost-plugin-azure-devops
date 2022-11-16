export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';

export const eventTypeBoards = {
    'workitem.created': 'Work item created',
    'workitem.updated': 'Work item updated',
    'workitem.deleted': 'Work item deleted',
    'workitem.commented': 'Work item commented on',
};

export const eventTypeRepos = {
    'git.pullrequest.created': 'Pull request created',
    'git.pullrequest.updated': 'Pull request updated',
    'ms.vss-code.git-pullrequest-comment-event': 'Pull request commented on',
    'git.pullrequest.merged': 'Pull request merge attempted',
    'git.push': 'Code Pushed',
};

export const eventTypePipelines = {
    'build.complete': 'Build completed',
    'ms.vss-release.release-abandoned-event': 'Release abandoned',
    'ms.vss-release.release-created-event': 'Release created',
    'ms.vss-release.deployment-approval-completed-event': 'Release deployment approval completed',
    'ms.vss-release.deployment-approval-pending-event': 'Release deployment approval pending',
    'ms.vss-release.deployment-completed-event': 'Release deployment completed',
    'ms.vss-release.deployment-started-event': 'Release deployment started',
    'ms.vss-pipelinechecks-events.approval-completed': 'Run stage approval completed',
    'ms.vss-pipelines.stage-state-changed-event': 'Run stage state changed',
    'ms.vss-pipelinechecks-events.approval-pending': 'Run stage waiting for approval',
    'ms.vss-pipelines.run-state-changed-event': 'Run state changed',
};

export const eventTypeMap: Record<EventType, string> = {
    ...eventTypeBoards,
    ...eventTypeRepos,
    ...eventTypePipelines,
};

export const boards = 'boards';
export const repos = 'repos';
export const pipelines = 'pipelines';

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
