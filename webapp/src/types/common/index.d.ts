/**
 * Keep all common types here which are to be used throughout the project
*/
type eventType = 'create' | 'update' | 'delete'
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
    label: string | JSX.Element;
    value: string;
}

type HttpMethod = 'GET' | 'POST';

type ApiServiceName = 'createTask' | 'testGet' | 'createLink' | 'getAllLinkedProjectsList' | 'unlinkProject' | 'getUserDetails' | 'getChannels' | 'createSubscription'

type PluginApiService = {
    path: string,
    method: HttpMethod,
    apiServiceName: ApiServiceName
}

type TabData = {
    title: string,
    tabPanel: JSX.Element
}

type CreateTaskFields = {
    title: string,
    description: string,
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

type APIRequestPayload = CreateTaskPayload | LinkPayload | ProjectDetails | UserDetails | ChannelList | {teamId: string} | void;

type DropdownOptionType = {
    label?: string | JSX.Element;
    value: string;
}

type FetchChannelParams = {
    teamId: string;
}

type FetchSubscriptionList = {
    project: string;
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

type SubscriptionDetails = {
    mattermostUserID: string
    projectID: string,
    projectName: string,
    organizationName: string,
    eventType: string,
    channelID: string,
    channelName: string,
}

type ModalId = 'linkProject' | 'createBoardTask' | null

type WebsocketEventParams = {
    event: string,
    data: Record<string, string>,
}
