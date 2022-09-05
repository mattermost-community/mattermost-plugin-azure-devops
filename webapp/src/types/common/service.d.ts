type HttpMethod = 'GET' | 'POST' | 'DELETE' ;

type ApiServiceName =
    'createTask' |
    'createLink' |
    'getAllLinkedProjectsList' |
    'unlinkProject' |
    'getUserDetails' |
    'createSubscription' |
    'getChannels' |
    'getSubscriptionList' |
    'deleteSubscription'

type PluginApiService = {
    path: string,
    method: HttpMethod,
    apiServiceName: ApiServiceName
}

type ApiErrorResponse = {
    data: {
        error: string
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
    void;
