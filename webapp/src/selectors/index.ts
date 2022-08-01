import plugin_constants from 'plugin_constants';

const pluginPrefix = `plugins-${plugin_constants.pluginId}`;

export const getTaskModalState = (state: any): TaskModalState => {
    return state[pluginPrefix].openTaskModalReducer;
};
