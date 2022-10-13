import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'pluginConstants';
import Utils from 'utils';

// Service to make plugin API requests
const azureDevOpsPluginApi = createApi({
    reducerPath: 'azureDevOpsPluginApi',
    baseQuery: fetchBaseQuery({baseUrl: Utils.getBaseUrls().pluginApiBaseUrl}),
    tagTypes: ['Posts'],
    endpoints: (builder) => ({
        [Constants.pluginApiServiceConfigs.createTask.apiServiceName]: builder.query<void, APIRequestData>({
            query: ({payload}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.createTask.path,
                method: Constants.pluginApiServiceConfigs.createTask.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createLink.apiServiceName]: builder.query<void, APIRequestData>({
            query: ({payload}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.createLink.path,
                method: Constants.pluginApiServiceConfigs.createLink.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName]: builder.query<ProjectDetails[], void>({
            query: () => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.path,
                method: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.unlinkProject.apiServiceName]: builder.query<void, APIRequestData>({
            query: ({payload}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.unlinkProject.path,
                method: Constants.pluginApiServiceConfigs.unlinkProject.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getUserDetails.apiServiceName]: builder.query<UserDetails, void>({
            query: () => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.getUserDetails.path,
                method: Constants.pluginApiServiceConfigs.getUserDetails.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createSubscription.apiServiceName]: builder.query<void, APIRequestData>({
            query: ({payload}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.createSubscription.path,
                method: Constants.pluginApiServiceConfigs.createSubscription.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getChannels.apiServiceName]: builder.query<ChannelList[], APIRequestData>({
            query: ({url}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Utils.addPathParamsToApiUrl(Constants.pluginApiServiceConfigs.getChannels.path, url?.pathParams),
                method: Constants.pluginApiServiceConfigs.getChannels.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName]: builder.query<SubscriptionDetails[], APIRequestData>({
            query: ({url}) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Utils.addPathParamsToApiUrl(Constants.pluginApiServiceConfigs.getSubscriptionList.path, url?.pathParams),
                method: Constants.pluginApiServiceConfigs.getSubscriptionList.method,
                params: {...url?.queryParams},
            }),
        }),
        [Constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName]: builder.query<void, APIRequestData>({
            query: (payload) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.deleteSubscription.path,
                method: Constants.pluginApiServiceConfigs.deleteSubscription.method,
                body: payload,
            }),
        }),
    }),
});

export default azureDevOpsPluginApi;
