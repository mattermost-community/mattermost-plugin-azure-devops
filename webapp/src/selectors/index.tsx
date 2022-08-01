import plugin_constants from 'plugin_constants';

const pluginPrefix = `plugins-${plugin_constants.pluginId}`;

// TODO: create a type for global state

export const getprojectDetailsState = (state: any) => {
    return state[pluginPrefix].projectDetailsSlice;
};

export const getLinkModalState = (state: any): LinkProjectModalState => {
    return state[pluginPrefix].openLinkModalReducer;
};

export const getRhsState = (state: any): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};

export const getTaskModalState = (state: any): TaskModalState => {
    return state[pluginPrefix].openTaskModalReducer;
};
