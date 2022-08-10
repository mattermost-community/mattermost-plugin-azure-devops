type LinkProjectModalState = {
    visibility: boolean,
    organization: string,
    project: string,
    isLinked: boolean,
}

type TaskModalState = {
    visibility: boolean,
}

type SubscribeModalState = {
    visibility: boolean,
    isLinked: boolean,
}

type userConnectionState = {
    isConnectionTriggered: boolean
    isUserDisconnected: boolean
}
