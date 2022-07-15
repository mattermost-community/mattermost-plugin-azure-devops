/**
 * Keep all plugin related constants here
*/

// Plugin configs
const pluginId = 'mattermost-plugin-azure-devops';

const AzureDevops = 'Azure Devops';
const RightSidebarHeader = 'Azure Devops';

// Plugin api service (RTK query) configs
const pluginApiServiceConfigs: Record<ApiServiceName, PluginApiService> = {
    fetchWellsList: {
        path: '/wells',
        method: 'GET',
        apiServiceName: 'fetchWellsList',
    },
    fetchWell: {
        path: '/wells',
        method: 'GET',
        apiServiceName: 'fetchWell',
    },
};

export default {
    pluginId,
    pluginApiServiceConfigs,
    AzureDevops,
    RightSidebarHeader,
};
