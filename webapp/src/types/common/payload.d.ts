type LinkPayload = {
    organization: string,
    project: string,
}

type CreateTaskPayload = {
    organization: string,
    project: string,
    type: string,
    fields: CreateTaskFields,
    timestamp: string
}

type SubscriptionPayload = {
    organization: string,
    project: string,
    eventType: string,
    serviceType: string,
    channelID: string,
    mmUserID: string,
    repository: string,
    targetBranch: string,
    repositoryName: string
    pullRequestCreatedBy: string
    pullRequestReviewersContains: string
    pullRequestCreatedByName: string
    pullRequestReviewersContainsName: string
    pushedBy: string
    pushedByName: string
    mergeResult: string
    mergeResultName: string
    notificationType: string
    notificationTypeName: string
    areaPath: string
}

interface PaginationQueryParams {
    page: number;
    per_page: number;
}

type ReposSubscriptionFiltersRequest = {
    organization: string
    project: string
    repository?: string
}

type SubscriptionFiltersPossibleValues = {
    displayValue: string
    value: string
}

type GetSubscriptionFiltersResponse = Record<string, SubscriptionFiltersPossibleValues[]>

type GetSubscriptionFiltersRequest = {
    organization: string
    projectId: string
    eventType: string
    filters: string[]
    repositoryId?: string
    releasePipelineId?: string
}
