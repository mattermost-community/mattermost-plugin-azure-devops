type LinkPayload = {
    organization: string,
    project: string,
}

type CreateTaskPayload = {
    organization: string,
    project: string,
    type: string,
    fields: CreateTaskFields,
}

type SubscriptionPayload = {
    organization: string,
    project: string,
    eventType: string,
    channelID: string
}
