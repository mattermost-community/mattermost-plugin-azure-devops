/**
 * Keep all common types here which are to be used throughout the project
*/

type HttpMethod = 'GET' | 'POST';

type ApiServiceName = 'createTask' | 'testGet' | 'createLink' | 'getAllLinkedProjectsList' | 'unlinkProject' | 'getUserDetails' | 'createSubscription' | 'getChannels'

type PluginApiService = {
    path: string,
    method: HttpMethod,
    apiServiceName: ApiServiceName
}

type PluginState = {
    'plugins-mattermost-plugin-azure-devops': RootState<{ [x: string]: QueryDefinition<void, BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, {}, FetchBaseQueryMeta>, never, WellList[], 'pluginApi'>; }, never, 'pluginApi'>
}

type TabData = {
    title: string,
    tabPanel: JSX.Element
}

type CreateTaskFields = {
    title: string,
    description: string,
}

type LinkPayload = {
    organization: string,
    project: string,
}

type CreateTaskPayload = {
    organization: string,
    project: string,
    type: string,
    fields: CreateTaskFields,
}

type SubscriptionPayload = {
    organization: string,
    project: string,
    eventType: string,
    channelID: string
}

type APIRequestPayload = CreateTaskPayload | LinkPayload | ProjectDetails | UserDetails | SubscriptionPayload | FetchChannelParams | void;

type DropdownOptionType = {
    label?: string | JSX.Element;
    value: string;
}

type TabsData = {
    title: string
    component: JSX.Element
}

type ProjectDetails = {
    mattermostID: string
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

type eventType = 'create' | 'update' | 'delete'

type SubscriptionDetails = {
    id: string
    name: string
    eventType: eventType
}

type ModalId = 'linkProject' | 'createBoardTask' | 'subscribeProject' | null
