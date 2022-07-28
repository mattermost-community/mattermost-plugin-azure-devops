import plugin_constants from 'plugin_constants';

const pluginPrefix = `plugins-${plugin_constants.pluginId}`;

export const getprojectDetailsState = (state: any) => {
    return state[pluginPrefix].projectDetailsSlice;
};

export const getRhsState = (state: any): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};
