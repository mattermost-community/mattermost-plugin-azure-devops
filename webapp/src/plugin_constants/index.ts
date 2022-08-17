/**
 * Keep all plugin related constants here
*/

// Plugin configs
const pluginId = 'mattermost-plugin-azure-devops';

const AzureDevops = 'Azure Devops';
const RightSidebarHeader = 'Azure Devops';

const MMCSRF = 'MMCSRF';
const HeaderCSRFToken = 'X-CSRF-Token';

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
        path: '/project/link',
        method: 'GET',
        apiServiceName: 'getAllLinkedProjectsList',
    },
    unlinkProject: {
        path: '/project/unlink',
        method: 'POST',
        apiServiceName: 'unlinkProject',
    },
    getUserDetails: {
        path: '/user',
        method: 'GET',
        apiServiceName: 'getUserDetails',
    },
    createSubscription: {
        path: '/subscriptions',
        method: 'POST',
        apiServiceName: 'createSubscription',
    },
    getChannels: {
        path: '/channels',
        method: 'GET',
        apiServiceName: 'getChannels',
    },
    getSubscriptionList: {
        path: '/subscriptions?project=',
        method: 'GET',
        apiServiceName: 'getSubscriptionList',
    },
    deleteSubscription: {
        path: '/subscriptions',
        method: 'DELETE',
        apiServiceName: 'deleteSubscription',
    },
};

export default {
    MMCSRF,
    HeaderCSRFToken,
    pluginId,
    pluginApiServiceConfigs,
    AzureDevops,
    RightSidebarHeader,
};
