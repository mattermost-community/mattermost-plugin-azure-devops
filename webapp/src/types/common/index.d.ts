/**
 * Keep all common types here which are to be used throughout the project
*/
// TODO: create enums for these types
type EventTypeBoards = 'workitem.created' | 'workitem.updated' | 'workitem.deleted' | 'workitem.commented'
type EventTypeRepos = 'git.pullrequest.created'| 'git.pullrequest.updated' | 'ms.vss-code.git-pullrequest-comment-event' | 'git.push' | 'git.pullrequest.merged'
type EventTypePipelines = 'build.complete' | 'ms.vss-release.release-abandoned-event' | 'ms.vss-release.release-created-event' | 'ms.vss-release.deployment-approval-completed-event' | 'ms.vss-release.deployment-approval-pending-event' | 'ms.vss-release.deployment-completed-event' | 'ms.vss-release.deployment-started-event' | 'ms.vss-pipelinechecks-events.approval-completed' | 'ms.vss-pipelines.stage-state-changed-event' | 'ms.vss-pipelinechecks-events.approval-pending' | 'ms.vss-pipelines.run-state-changed-event'
type EventType = EventTypeBoards | EventTypeRepos | EventTypePipelines
type ModalId = 'linkProject' | 'createBoardTask' | 'subscribeProject' | null

type TabData = {
    title: string,
    tabPanel: JSX.Element
}

type TabsData = {
    title: string
    component: JSX.Element
}

type LabelValuePair = {
    label?: string | JSX.Element;
    value: string;
    metaData?: string;
}

type CreateTaskFields = {
    title: string,
    description: string,
    areaPath: string,
}

type ProjectDetails = {
    mattermostUserID: string
    projectID: string,
    projectName: string,
    organizationName: string
}

type UserDetails = {
    MattermostUserID: string
}

type ChannelList = {
    display_name: string,
    id: string,
    name: string,
    team_id: string,
    team_name: string,
    type: string
}

type FetchChannelParams = {
    teamId: string;
}

interface FetchSubscriptionList extends PaginationQueryParams {
    project: string;
    channel_id: string;
    created_by: string;
    team_id: string;
    service_type: string;
    event_type: string;
}

type SubscriptionDetails = {
    mattermostUserID: string
    projectID: string,
    projectName: string,
    organizationName: string,
    eventType: string,
    serviceType: string,
    channelID: string,
    channelName: string,
    channelType: string,
    createdBy: string,
    repository: string,
    targetBranch: string,
    repositoryName: string
    pullRequestCreatedBy: string
    pullRequestReviewersContains: string
    pullRequestCreatedByName: string
    pullRequestReviewersContainsName: string
    pushedBy: string
    pushedByName: string
    mergeResult: string
    mergeResultName: string
    notificationType: string
    notificationTypeName: string
    areaPath: string
    buildPipeline: string
    buildStatus: string
    releasePipeline: string
    stageName: string
    approvalType: string
    approvalStatus: string
    releaseStatus: string
    buildStatusName: string
    releasePipelineName: string
    stageNameValue: string
    approvalTypeName: string
    approvalStatusName: string
    releaseStatusName: string
    runPipeline: string
    runStage: string
    runEnvironment: string
    runStageId: string
    runStageStateId: string
    runStageResultId: string
    runStateId: string
    runResultId: string
}

type WebsocketEventParams = {
    event: string,
    data: Record<string, string>,
}

type ConfirmationModalErrorPanelProps = {
    title: string,
    onSecondaryBtnClick: () => void,
}

type SubscriptionFilters = {
    createdBy: string,
    serviceType: string,
    eventType: string,
}
