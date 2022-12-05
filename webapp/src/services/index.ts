import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'pluginConstants';
import Utils, {addPathParamsToApiUrl} from 'utils';

// Service to make plugin API requests
const azureDevOpsPluginApi = createApi({
    reducerPath: 'azureDevOpsPluginApi',
    baseQuery: fetchBaseQuery({baseUrl: Utils.getBaseUrls().pluginApiBaseUrl}),
    tagTypes: ['Posts'],
    endpoints: (builder) => ({
        [Constants.pluginApiServiceConfigs.createTask.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.createTask.path,
                method: Constants.pluginApiServiceConfigs.createTask.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createLink.apiServiceName]: builder.query<CreateLinkResponse, APIRequestPayload>({
            query: (payload) => ({
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
        [Constants.pluginApiServiceConfigs.unlinkProject.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
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
        [Constants.pluginApiServiceConfigs.createSubscription.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.createSubscription.path,
                method: Constants.pluginApiServiceConfigs.createSubscription.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getChannels.apiServiceName]: builder.query<ChannelList[], FetchChannelParams>({
            query: (params) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: `${Constants.pluginApiServiceConfigs.getChannels.path}/${params.teamId}`,
                method: Constants.pluginApiServiceConfigs.getChannels.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName]: builder.query<SubscriptionDetails[], FetchSubscriptionList>({
            query: (params) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: `${Constants.pluginApiServiceConfigs.getSubscriptionList.path}/${params.team_id}`,
                method: Constants.pluginApiServiceConfigs.getSubscriptionList.method,
                params: {...params},
            }),
        }),
        [Constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.deleteSubscription.path,
                method: Constants.pluginApiServiceConfigs.deleteSubscription.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName]: builder.query<GetSubscriptionFiltersResponse, GetSubscriptionFiltersRequest>({
            query: (payload) => ({
                headers: {[Constants.common.HeaderCSRFToken]: Cookies.get(Constants.common.MMCSRF)},
                url: Constants.pluginApiServiceConfigs.getSubscriptionFilters.path,
                method: Constants.pluginApiServiceConfigs.getSubscriptionFilters.method,
                body: payload,
            }),
        }),
    }),
});

export default azureDevOpsPluginApi;
