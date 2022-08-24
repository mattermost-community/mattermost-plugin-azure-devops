// Create subscription modal
const eventTypeOptions: LabelValuePair[] = [
    {
        value: 'create',
        label: 'Create',
    },
    {
        value: 'update',
        label: 'Update',
    },
    {
        value: 'delete',
        label: 'Delete',
    },
];

export const subscriptionModal: Record<SubscriptionModalFields, ModalFormFieldConfig> = {
    organization: {
        label: 'Organization name',
        type: 'dropdown',
        value: '',
        validations: {
            isRequired: true,
        },
    },
    project: {
        label: 'Project name',
        value: '',
        type: 'dropdown',
        validations: {
            isRequired: true,
        },
    },
    eventType: {
        label: 'Event type',
        value: '',
        type: 'dropdown',
        optionsList: eventTypeOptions,
        validations: {
            isRequired: true,
        },
    },
    channelID: {
        label: 'Channel name',
        value: '',
        type: 'dropdown',
        validations: {
            isRequired: true,
        },
    },

    // add 'timestamp' field only if you don't want to use cached RTK Api query
    timestamp: {
        label: 'time',
        type: 'timestamp',
        value: '',
    },
};

// Create task modal
const taskTypeOptions = [
    {
        value: 'Task',
        label: 'Task',
    },
    {
        value: 'Epic',
        label: 'Epic',
    },
    {
        value: 'Issue',
        label: 'Issue',
    },
];

export const createTaskModal: Record<CreateTaskModalFields, ModalFormFieldConfig> = {
    organization: {
        label: 'Organization name',
        type: 'dropdown',
        value: '',
        validations: {
            isRequired: true,
        },
    },
    project: {
        label: 'Project name',
        value: '',
        type: 'dropdown',
        validations: {
            isRequired: true,
        },
    },
    type: {
        label: 'Work item type',
        value: '',
        type: 'dropdown',
        optionsList: taskTypeOptions,
        validations: {
            isRequired: true,
        },
    },
    title: {
        label: 'Title',
        value: '',
        type: 'text',
        validations: {
            isRequired: true,
        },
    },
    description: {
        label: 'Description',
        value: '',
        type: 'text',
    },

    // add 'timestamp' field only if you don't want to use cached RTK Api query
    timestamp: {
        label: 'time',
        type: 'timestamp',
        value: '',
    },
};
