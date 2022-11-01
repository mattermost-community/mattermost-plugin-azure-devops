export const pluginId = 'mattermost-plugin-azure-devops';

export const AzureDevops = 'Azure DevOps';
export const RightSidebarHeader = 'Azure DevOps';

export const MMCSRF = 'MMCSRF';
export const HeaderCSRFToken = 'X-CSRF-Token';

export const eventTypeBoards = {
    'workitem.created': 'Work Item Created',
    'workitem.updated': 'Work Item Updated',
    'workitem.deleted': 'Work Item Deleted',
    'workitem.commented': 'Work Item Commented On',
};

export const eventTypeRepos = {
    'git.pullrequest.created': 'Pull Request Created',
    'git.pullrequest.updated': 'Pull Request Updated',
    'ms.vss-code.git-pullrequest-comment-event': 'Pull Request Commented On',
    'git.pullrequest.merged': 'Pull Request Merge Attempted',
    'git.push': 'Code Pushed',
};

export const eventTypePipelines = {
    'build.complete': 'Build Completed',
    'ms.vss-release.release-abandoned-event': 'Release Abandoned',
    'ms.vss-release.release-created-event': 'Release Created',
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
    'build.complete': 'Build Completed',
    'ms.vss-release.release-abandoned-event': 'Release Abandoned',
    'ms.vss-release.release-created-event': 'Release Created',
    'ms.vss-release.deployment-approval-completed-event': 'Release deployment approval completed',
    'ms.vss-release.deployment-approval-pending-event': 'Release deployment approval pending',
    'ms.vss-release.deployment-completed-event': 'Release deployment completed',
    'ms.vss-release.deployment-started-event': 'Release deployment started',
    'ms.vss-pipelinechecks-events.approval-completed': 'Run stage approval completed',
    'ms.vss-pipelines.stage-state-changed-event': 'Run stage state changed',
    'ms.vss-pipelinechecks-events.approval-pending': 'Run stage waiting for approval',
    'ms.vss-pipelines.run-state-changed-event': 'Run state changed',
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
