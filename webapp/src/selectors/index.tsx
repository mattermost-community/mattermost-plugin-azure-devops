export const getProjectDetailsState = (state: ReduxState): ProjectDetails => {
    return state['plugins-mattermost-plugin-azure-devops'].projectDetailsSlice;
};

export const getRhsState = (state: ReduxState): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};

export const getGlobalModalState = (state: ReduxState): GlobalModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].globalModalSlice;
};

export const getLinkModalState = (state: ReduxState): LinkProjectModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].linkProjectModalSlice;
};

export const getCreateTaskModalState = (state: ReduxState): CreateTaskModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].createTaskModalSlice;
};

export const getSubscribeModalState = (state: ReduxState): SubscribeModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].subscriptionModalSlice;
};

export const getApiRequestCompletionState = (state: ReduxState): ApiRequestCompletionState => {
    return state['plugins-mattermost-plugin-azure-devops'].apiRequestCompletionSlice;
};

export const getApiQueriesState = (state: ReduxState): ApiQueriesState => {
    return state['plugins-mattermost-plugin-azure-devops'].azureDevOpsPluginApi?.queries;
};

export const getWebsocketEventState = (state: ReduxState): WebsocketEventState => {
    return state['plugins-mattermost-plugin-azure-devops'].websocketEventSlice;
};
