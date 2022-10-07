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
}

interface PaginationQueryParams {
    page: number;
    per_page: number;
}
