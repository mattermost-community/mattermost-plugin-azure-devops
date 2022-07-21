/**
 * Keep all plugin related constants here
*/

// Plugin configs
const pluginId = 'mattermost-plugin-azure-devops';

const AzureDevops = 'Azure Devops';
const RightSidebarHeader = 'Azure Devops';

const MMUSERID = 'MMUSERID';
const HeaderMattermostUserID = 'User-ID';

// Plugin api service (RTK query) configs
const pluginApiServiceConfigs: Record<ApiServiceName, PluginApiService> = {
    createTask: {
        path: '/tasks',
        method: 'POST',
        apiServiceName: 'createTask',
    },
    testGet: {
        path: '/test',
        method: 'GET',
        apiServiceName: 'testGet',
    },
};

export default {
    MMUSERID,
    HeaderMattermostUserID,
    pluginId,
    pluginApiServiceConfigs,
    AzureDevops,
    RightSidebarHeader,
};
