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
    'deleteSubscription' |
    'getRepositories' |
    'getRepositoryBranches' |
    'getSubscriptionFilters'

type PluginApiService = {
    path: string,
    method: HttpMethod,
    apiServiceName: ApiServiceName
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
    ReposSubscriptionFiltersRequest |
    ReposSubscriptionFiltersResponse |
    ReposSubscriptionTargetBranchFilterResponse |
    GetSubscriptionFiltersRequest |
    GetSubscriptionFiltersResponse |
    void;
