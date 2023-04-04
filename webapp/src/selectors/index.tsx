export const getRhsState = (state: ReduxState): {isSidebarOpen: boolean} => state.views.rhs;

// TODO: create constants for string literals to prevent changing hard coded string values at multiple places in case of modifications
// Plugin state
const getPluginState = (state: ReduxState) => state['plugins-mattermost-plugin-azure-devops'];

export const getProjectDetailsState = (state: ReduxState): ProjectDetails => getPluginState(state).projectDetailsSlice;

export const getGlobalModalState = (state: ReduxState): GlobalModalState => getPluginState(state).globalModalSlice;

export const getLinkModalState = (state: ReduxState): LinkProjectModalState => getPluginState(state).linkProjectModalSlice;

export const getCreateTaskModalState = (state: ReduxState): CreateTaskModalState => getPluginState(state).createTaskModalSlice;

export const getSubscribeModalState = (state: ReduxState): SubscribeModalState => getPluginState(state).subscriptionModalSlice;

export const getApiRequestCompletionState = (state: ReduxState): ApiRequestCompletionState => getPluginState(state).apiRequestCompletionSlice;

export const getApiQueriesState = (state: ReduxState): ApiQueriesState => getPluginState(state).azureDevOpsPluginApi?.queries;

export const getWebsocketEventState = (state: ReduxState): WebsocketEventState => getPluginState(state).websocketEventSlice;
