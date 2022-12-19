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
    buildPipeline: string
    buildStatus: string
    releasePipeline: string
    stageName: string
    approvalType: string
    approvalStatus: string
    releaseStatus: string
    buildStatusName: string
    releasePipelineName: string
    stageNameValue: string
    approvalTypeName: string
    approvalStatusName: string
    releaseStatusName: string
}

interface PaginationQueryParams {
    page: number;
    per_page: number;
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
