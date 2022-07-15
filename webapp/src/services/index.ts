import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Constants from 'plugin_constants';
import Utils from 'utils';

// Service to make plugin API requests
const pluginApi = createApi({
    reducerPath: 'pluginApi',
    baseQuery: fetchBaseQuery({baseUrl: Utils.getBaseUrls().pluginApiBaseUrl}),
    endpoints: (builder) => ({
        [Constants.pluginApiServiceConfigs.fetchWellsList.apiServiceName]: builder.query<WellList[], void>({
            query: () => ({
                url: Constants.pluginApiServiceConfigs.fetchWellsList.path,
                method: Constants.pluginApiServiceConfigs.fetchWellsList.method,
            }),
        }),
        [Constants.pluginApiServiceConfigs.fetchWell.apiServiceName]: builder.query<WellList[], void>({
            query: () => ({
                url: Constants.pluginApiServiceConfigs.fetchWell.path,
                method: Constants.pluginApiServiceConfigs.fetchWell.method,
            }),
        }),
    }),
});

export default pluginApi;
