/**
 * Keep all common types here which are to be used throughout the project
*/

type HttpMethod = 'GET' | 'POST';

type ApiServiceName = 'fetchWellsList' | 'fetchWell'

type PluginApiService = {
    path: string,
    method: httpMethod,
    apiServiceName: string
}

type PluginState = {
    'plugins-mattermost-plugin-wellsite-witsml': RootState<{ [x: string]: QueryDefinition<void, BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, {}, FetchBaseQueryMeta>, never, WellList[], 'pluginApi'>; }, never, 'pluginApi'>
}

type TabData = {
    title: string,
    tabPanel: JSX.Element
}
