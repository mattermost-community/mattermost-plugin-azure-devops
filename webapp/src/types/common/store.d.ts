interface ReduxState extends GlobalState {
    views: {
        rhs: {
            isSidebarOpen: boolean
        }
    }
    'plugins-mattermost-plugin-azure-devops': RootState<{ [x: string]: QueryDefinition<void, BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, {}, FetchBaseQueryMeta>, never, void, 'azureDevOpsPluginApi'>; }, never, 'azureDevOpsPluginApi'>
}

type GlobalModalState = {
    modalId: ModalId
    commandArgs: Array<string>
}

type GlobalModalActionPayload = {
    isVisible: boolean
    commandArgs: Array<string>
    args?: Array<string>
    isActionDone?: boolean
}

type LinkProjectModalState = {
    visibility: boolean,
    organization: string,
    project: string,
    isLinked: boolean,
}

type SubscribeModalState = {
    visibility: boolean,
    isCreated: boolean,
    serviceType: string,
    organization?: string | null,
    project?: string | null,
    projectID?: string | null
}

type CreateTaskCommandArgs = {
    title: string;
    description: string;
}

type CreateTaskModalState = {
    visibility: boolean
    commandArgs: CreateTaskCommandArgs
}

type ApiQueriesState = {
    [key: string]: Record<string, string>
}

type ApiRequestCompletionState = {
    requests: PluginApiServiceName[]
}

type ProjectDetails = {
    mattermostUserID: string
    projectID: string,
    projectName: string,
    organizationName: string
}

type WebsocketEventState = {
    isConnected: boolean;
    isSubscriptionDeleted: boolean;
};
