// Plugin api service (RTK query) configs
export const pluginApiServiceConfigs: Record<ApiServiceName, PluginApiService> = {
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
        path: '/subscriptions',
        method: 'GET',
        apiServiceName: 'getSubscriptionList',
    },
    deleteSubscription: {
        path: '/subscriptions',
        method: 'DELETE',
        apiServiceName: 'deleteSubscription',
    },
};
