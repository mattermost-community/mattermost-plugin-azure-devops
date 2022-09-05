type GlobalModalState = {
    modalId: ModalId
    commandArgs: Array<string>
}

type GlobalModalActionPayload = {
    isVisible: boolean
    commandArgs: Array<string>
}

type LinkProjectModalState = {
    visibility: boolean,
    organization: string,
    project: string,
    isLinked: boolean,
}

type CreateTaskCommandArgs = {
    title: string;
    description: string;
}

type CreateTaskModalState = {
    visibility: boolean
    commandArgs: CreateTaskCommandArgs
}
