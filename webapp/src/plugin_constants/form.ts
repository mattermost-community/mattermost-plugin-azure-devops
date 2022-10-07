import {subscriptionFilters} from './common';

// Create subscription modal
export const boardEventTypeOptions: LabelValuePair[] = [
    {
        value: 'workitem.created',
        label: 'Create',
    },
    {
        value: 'workitem.updated',
        label: 'Update',
    },
    {
        value: 'workitem.deleted',
        label: 'Delete',
    },
    {
        value: 'workitem.commented',
        label: 'Comment',
    },
];

export const repoEventTypeOptions: LabelValuePair[] = [
    {
        value: 'git.pullrequest.created',
        label: 'Create',
    },
    {
        value: 'git.pullrequest.updated',
        label: 'Update',
    },
    {
        value: 'ms.vss-code.git-pullrequest-comment-event',
        label: 'Comment',
    },
    {
        value: 'git.push',
        label: 'Code Push',
    },
    {
        value: 'git.pullrequest.merged',
        label: 'Merge Attempt',
    },
];

const serviceTypeOptions: LabelValuePair[] = [
    {
        value: 'boards',
        label: 'Board',
    },
    {
        value: 'repos',
        label: 'Repos',
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
    serviceType: {
        label: 'Service type',
        value: 'boards',
        type: 'dropdown',
        optionsList: serviceTypeOptions,
        validations: {
            isRequired: true,
        },
    },
    eventType: {
        label: 'Event type',
        value: '',
        type: 'dropdown',
        optionsList: boardEventTypeOptions,
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

    // add 'timestamp' field only if you don't want to use cached RTK API query
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
    areaPath: {
        label: 'Area Path',
        value: '',
        type: 'text',
    },

    // add 'timestamp' field only if you don't want to use cached RTK API query
    timestamp: {
        label: 'time',
        type: 'timestamp',
        value: '',
    },
};

// Link project modal
export const linkProjectModal: Record<LinkProjectModalFields, ModalFormFieldConfig> = {
    organization: {
        label: 'Organization name',
        type: 'text',
        value: '',
        validations: {
            isRequired: true,
        },
    },
    project: {
        label: 'Project name',
        value: '',
        type: 'text',
        validations: {
            isRequired: true,
        },
    },

    // add 'timestamp' field only if you don't want to use cached RTK API query
    timestamp: {
        label: 'time',
        type: 'timestamp',
        value: '',
    },
};

export const subscriptionFilterCreatedByOptions = [
    {
        value: subscriptionFilters.createdBy.me,
        label: 'Me',
    },
    {
        value: subscriptionFilters.createdBy.anyone,
        label: 'Anyone',
    },
];

export const subscriptionFilterServiceTypeOptions = [
    {
        value: subscriptionFilters.serviceType.boards,
        label: 'Boards',
    },
    {
        value: subscriptionFilters.serviceType.repo,
        label: 'Repos',
    },
];
