import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'pluginConstants';
import Utils from 'utils';

// Service to make plugin API requests
export const azureDevOpsPluginApi = createApi({
    reducerPath: 'azureDevOpsPluginApi',
    baseQuery: fetchBaseQuery({
        baseUrl: Utils.getBaseUrls().pluginApiBaseUrl,
        prepareHeaders: (headers) => {
            headers.set(Constants.common.HeaderCSRFToken, Cookies.get(Constants.common.MMCSRF) ?? '');

            return headers;
        },
    }),
    tagTypes: ['Posts'],
    endpoints: (builder) => ({
        [Constants.pluginApiServiceConfigs.createTask.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.createTask.path,
                method: Constants.pluginApiServiceConfigs.createTask.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createLink.apiServiceName]: builder.query<CreateLinkResponse, APIRequestPayload>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.createLink.path,
                method: Constants.pluginApiServiceConfigs.createLink.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName]: builder.query<ProjectDetails[], void>({
            query: () => ({
                url: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.path,
                method: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.unlinkProject.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.unlinkProject.path,
                method: Constants.pluginApiServiceConfigs.unlinkProject.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getUserDetails.apiServiceName]: builder.query<UserDetails, void>({
            query: () => ({
                url: Constants.pluginApiServiceConfigs.getUserDetails.path,
                method: Constants.pluginApiServiceConfigs.getUserDetails.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createSubscription.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.createSubscription.path,
                method: Constants.pluginApiServiceConfigs.createSubscription.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName]: builder.query<SubscriptionDetails[], FetchSubscriptionList>({
            query: (params) => ({
                url: `${Constants.pluginApiServiceConfigs.getSubscriptionList.path}/${params.team_id}/${params.organization}/${params.project}`,
                method: Constants.pluginApiServiceConfigs.getSubscriptionList.method,
                params: {...params},
            }),
        }),
        [Constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.deleteSubscription.path,
                method: Constants.pluginApiServiceConfigs.deleteSubscription.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName]: builder.query<GetSubscriptionFiltersResponse, GetSubscriptionFiltersRequest>({
            query: (payload) => ({
                url: Constants.pluginApiServiceConfigs.getSubscriptionFilters.path,
                method: Constants.pluginApiServiceConfigs.getSubscriptionFilters.method,
                body: payload,
            }),
        }),
    }),
});
