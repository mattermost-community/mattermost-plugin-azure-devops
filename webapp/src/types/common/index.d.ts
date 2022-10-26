/**
 * Keep all common types here which are to be used throughout the project
*/
// TODO: create enums for these types
type EventTypeBoards = 'workitem.created' | 'workitem.updated' | 'workitem.deleted' | 'workitem.commented'
type EventTypeRepos = 'git.pullrequest.created'| 'git.pullrequest.updated' | 'ms.vss-code.git-pullrequest-comment-event' | 'git.push' | 'git.pullrequest.merged'
type EventType = EventTypeBoards | EventTypeRepos
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

type ReposSubscriptionFiltersResponse = {
    id: string,
    name: string
}
