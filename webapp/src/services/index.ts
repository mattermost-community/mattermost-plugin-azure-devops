import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'plugin_constants';
import Utils from 'utils';

// Service to make plugin API requests
const pluginApi = createApi({
  reducerPath: 'pluginApi',
  baseQuery: fetchBaseQuery({ baseUrl: Utils.getBaseUrls().pluginApiBaseUrl }),
  tagTypes: ['Posts'],

  endpoints: (builder) => ({
    // TODO: example of GET request, remove later if not required
    [Constants.pluginApiServiceConfigs.testGet.apiServiceName]: builder.query<any, void>({
      query: () => ({
        headers: { [Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID) },
        url: Constants.pluginApiServiceConfigs.testGet.path,
        method: Constants.pluginApiServiceConfigs.testGet.method,
      }),
    }),
    // TODO: example of POST request, remove later if not required
    [Constants.pluginApiServiceConfigs.createTask.apiServiceName]: builder.query<void, CreateTaskPayload>({
      query: (payload) => ({
        headers: { [Constants.HeaderMattermostUserID]: Cookies.get(Constants.MMUSERID) },
        url: Constants.pluginApiServiceConfigs.createTask.path,
        method: Constants.pluginApiServiceConfigs.createTask.method,
        body: payload
      }),
    })
  }),
});

export default pluginApi;
