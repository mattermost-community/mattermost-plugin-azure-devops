type PluginState = {
    'plugins-mattermost-plugin-azure-devops': RootState<{ [x: string]: QueryDefinition<void, BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, {}, FetchBaseQueryMeta>, never, void, 'azureDevOpsPluginApi'>; }, never, 'pluginApi'>
}

type ApiRequestCompletionState = {
    requests: ApiServiceName[]
}

type GlobalModalState = {
    modalId: ModalId
    commandArgs: Array<string>
}

type GlobalModalActionPayload = {
    isVisible: boolean
    commandArgs: Array<string>
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
}

type TaskFieldsCommandArgs = {
    title: string;
    description: string;
}

type CreateTaskModalState = {
    visibility: boolean
    commandArgs: TaskFieldsCommandArgs
}

type ApiQueriesState = {
    [key: string]: Record<string, string>
}
