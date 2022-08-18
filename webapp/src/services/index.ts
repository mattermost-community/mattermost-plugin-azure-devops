import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'plugin_constants';
import Utils from 'utils';

// Service to make plugin API requests
const pluginApi = createApi({
    reducerPath: 'pluginApi',
    baseQuery: fetchBaseQuery({baseUrl: Utils.getBaseUrls().pluginApiBaseUrl}),
    tagTypes: ['Posts'],
    endpoints: (builder) => ({
        [Constants.pluginApiServiceConfigs.createTask.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID)},
                url: Constants.pluginApiServiceConfigs.createTask.path,
                method: Constants.pluginApiServiceConfigs.createTask.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.createLink.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID)},
                url: Constants.pluginApiServiceConfigs.createLink.path,
                method: Constants.pluginApiServiceConfigs.createLink.method,
                body: payload,
            }),
        }),
        [Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName]: builder.query<ProjectDetails[], void>({
            query: () => ({
                headers: {[Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID)},
                url: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.path,
                method: Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.unlinkProject.apiServiceName]: builder.query<void, APIRequestPayload>({
            query: (payload) => ({
                headers: {[Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID)},
                url: Constants.pluginApiServiceConfigs.unlinkProject.path,
                method: Constants.pluginApiServiceConfigs.unlinkProject.method,
                body: payload,
            }),
        }),
    }),
});

export default pluginApi;
