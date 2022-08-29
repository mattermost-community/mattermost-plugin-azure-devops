import plugin_constants from 'plugin_constants';

const pluginPrefix = `plugins-${plugin_constants.common.pluginId}`;

// TODO: create a type for global state

export const getGlobalModalState = (state: any): GlobalModalState => {
    return state[pluginPrefix].globalModalSlice;
};

export const getProjectDetailsState = (state: any) => {
    return state[pluginPrefix].projectDetailsSlice;
};

export const getLinkModalState = (state: any): LinkProjectModalState => {
    return state[pluginPrefix].linkProjectModalSlice;
};

export const getCreateTaskModalState = (state: any): CreateTaskModalState => {
    return state[pluginPrefix].openTaskModalReducer;
};

export const getRhsState = (state: any): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};

export const getSubscribeModalState = (state: any): SubscribeModalState => {
    return state[pluginPrefix].subscriptionModalSlice;
};

export const getApiRequestCompletionState = (state: any): ApiRequestCompletionState => {
    return state[pluginPrefix].apiRequestCompletionSlice;
};

export const getApiQueriesState = (state: any): ApiQueriesState => {
    return state[pluginPrefix].azureDevopsPluginApi?.queries;
};
