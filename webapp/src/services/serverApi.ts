import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react';

import Cookies from 'js-cookie';

import Constants from 'pluginConstants';
import Utils from 'utils';

// Service to make Mattermost's server API requests
export const mattermostServerApi = createApi({
    reducerPath: 'mattermostServerApi',
    baseQuery: fetchBaseQuery({
        baseUrl: Utils.getBaseUrls().mattermostApiBaseUrl,
        prepareHeaders: (headers) => {
            headers.set('authorization', `Bearer ${Cookies.get(Constants.common.MMAUTHTOKEN)}`);
            headers.set(Constants.common.HeaderCSRFToken, Cookies.get(Constants.common.MMCSRF) ?? '');

            return headers;
        },
    }),
    tagTypes: ['Posts'],
    endpoints: (builder) => ({
        [Constants.mattermostApiServiceConfigs.getChannels.apiServiceName]: builder.query<ChannelList[], FetchChannelParams>({
            query: (params) => {
                const currentUserId = Cookies.get(Constants.common.MMUSERID) ?? '';
                return ({
                    url: Constants.mattermostApiServiceConfigs.getChannels.path([currentUserId, params.teamId]),
                    method: Constants.mattermostApiServiceConfigs.getChannels.method,
                });
            },
        }),
    }),
});
