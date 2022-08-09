import plugin_constants from 'plugin_constants';

const pluginPrefix = `plugins-${plugin_constants.pluginId}`;

export const getprojectDetailsState = (state: any) => {
    return state[pluginPrefix].projectDetailsSlice;
};
