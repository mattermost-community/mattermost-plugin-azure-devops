type HttpMethod = 'GET' | 'POST' | 'DELETE' ;

type PluginApiServiceName =
    'createTask' |
    'createLink' |
    'getAllLinkedProjectsList' |
    'unlinkProject' |
    'getUserDetails' |
    'createSubscription' |
    'getSubscriptionList' |
    'deleteSubscription' |
    'getSubscriptionFilters'

type PluginApiService = {
    path: string,
    method: HttpMethod,
    apiServiceName: PluginApiServiceName
}

type ApiErrorResponse = {
    data: {
        Error: string
    },
    status: number
}

type APIRequestPayload =
    CreateTaskPayload |
    LinkPayload |
    ProjectDetails |
    UserDetails |
    SubscriptionPayload |
    FetchChannelParams |
    FetchSubscriptionList |
    GetSubscriptionFiltersRequest |
    GetSubscriptionFiltersResponse |
    void;
