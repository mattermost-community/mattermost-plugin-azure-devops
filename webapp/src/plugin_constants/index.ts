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
    createLink: {
        path: '/link',
        method: 'POST',
        apiServiceName: 'createLink',
    },
    testGet: {
        path: '/test',
        method: 'GET',
        apiServiceName: 'testGet',
    },
    getAllLinkedProjectsList: {
        path: '/link/project',
        method: 'GET',
        apiServiceName: 'getAllLinkedProjectsList'
    }
};

export default {
    MMUSERID,
    HeaderMattermostUserID,
    pluginId,
    pluginApiServiceConfigs,
    AzureDevops,
    RightSidebarHeader,
};
