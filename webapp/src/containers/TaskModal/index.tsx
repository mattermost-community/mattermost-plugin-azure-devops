import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';

import usePluginApi from 'hooks/usePluginApi';
import {toggleShowTaskModal} from 'reducers/taskModal';
import {getCreateTaskModalState} from 'selectors';

// TODO: fetch the organization and project options from API later.
const organizationOptions = [
    {
        value: 'bs-test',
        label: 'bs-test',
    },
    {
        value: 'brightscout-test',
        label: 'brightscout-test',
    },
];

const projectOptions = [
    {
        value: 'bs',
        label: 'bs',
    },
    {
        value: 'bs-2',
        label: 'bs-2',
    },
    {
        value: 'bs-3',
        label: 'bs-3',
    },
    {
        value: 'azure-test',
        label: 'azure-test',
    },
];

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

const TaskModal = () => {
    // State variables
    const [taskDetails, setTaskDetails] = useState<CreateTaskPayload>({
        organization: '',
        project: '',
        type: '',
        fields: {
            title: '',
            description: '',
        },
    });
    const [taskDetailsError, setTaskDetailsError] = useState<CreateTaskPayload>({
        organization: '',
        project: '',
        type: '',
        fields: {
            title: '',
            description: '',
        },
    });

    // Hooks
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    // Function to hide the modal and reset all the states.
    const resetModalState = () => {
        setTaskDetails({
            organization: '',
            project: '',
            type: '',
            fields: {
                title: '',
                description: '',
            },
        });
        setTaskDetailsError({
            organization: '',
            project: '',
            type: '',
            fields: {
                title: '',
                description: '',
            },
        });
        dispatch(toggleShowTaskModal({isVisible: false, commandArgs: []}));
    };

    const onOrganizationChange = (value: string) => {
        setTaskDetailsError({...taskDetailsError, organization: ''});
        setTaskDetails({...taskDetails, organization: value});
    };

    const onProjectChange = (value: string) => {
        setTaskDetailsError({...taskDetailsError, project: ''});
        setTaskDetails({...taskDetails, project: value});
    };

    const onTaskTypeChange = (value: string) => {
        setTaskDetailsError({...taskDetailsError, type: ''});
        setTaskDetails({...taskDetails, type: value});
    };

    const onTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTaskDetailsError({...taskDetailsError, fields: {...taskDetailsError.fields, title: ''}});
        setTaskDetails({...taskDetails, fields: {...taskDetails.fields, title: (e.target as HTMLInputElement).value}});
    };

    const onDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTaskDetailsError({...taskDetailsError, fields: {...taskDetailsError.fields, description: ''}});
        setTaskDetails({...taskDetails, fields: {...taskDetails.fields, description: (e.target as HTMLInputElement).value}});
    };

    const onConfirm = () => {
        const errorState: CreateTaskPayload = {
            organization: '',
            project: '',
            type: '',
            fields: {
                title: '',
                description: '',
            },
        };

        if (taskDetails.organization === '') {
            errorState.organization = 'Organization is required';
        }
        if (taskDetails.project === '') {
            errorState.project = 'Project is required';
        }
        if (taskDetails.type === '') {
            errorState.type = 'Work item type is required';
        }
        if (taskDetails.fields.title === '') {
            errorState.fields.title = 'Title is required';
        }

        if (errorState.organization || errorState.project || errorState.type || errorState.fields.title) {
            setTaskDetailsError(errorState);
            return;
        }

        // Make POST api request
        createTask();
    };

    // Make POST api request to create a task
    const createTask = async () => {
        const createTaskResponse = await usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails);
        if (createTaskResponse) {
            resetModalState();
        }
    };

    // Set modal field values
    useEffect(() => {
        if (getCreateTaskModalState(usePlugin.state).visibility) {
            // Pre-select the dropdown value in case of single option.
            if (organizationOptions.length === 1) {
                setTaskDetails({...taskDetails, organization: organizationOptions[0].value});
            }
            if (projectOptions.length === 1) {
                setTaskDetails({...taskDetails, project: projectOptions[0].value});
            }

            setTaskDetails({
                ...taskDetails,
                fields: {
                    title: getCreateTaskModalState(usePlugin.state).commandArgs.title,
                    description: getCreateTaskModalState(usePlugin.state).commandArgs.description,
                },
            });
        }
    }, [getCreateTaskModalState(usePlugin.state).visibility]);

    const apiResponse = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails);
    return (
        <Modal
            show={getCreateTaskModalState(usePlugin.state).visibility}
            title='Create Task'
            onHide={resetModalState}
            onConfirm={onConfirm}
            confirmBtnText='Create task'
            loading={apiResponse.isLoading}
            confirmDisabled={apiResponse.isLoading}
        >
            <>
                <Dropdown
                    placeholder='Organization name'
                    value={taskDetails.organization}
                    onChange={onOrganizationChange}
                    options={organizationOptions}
                    required={true}
                    error={taskDetailsError.organization}
                    disabled={apiResponse.isLoading}
                />
                <Dropdown
                    placeholder='Project name'
                    value={taskDetails.project}
                    onChange={onProjectChange}
                    options={projectOptions}
                    required={true}
                    error={taskDetailsError.project}
                    disabled={apiResponse.isLoading}
                />
                <Dropdown
                    placeholder='Work item type'
                    value={taskDetails.type}
                    onChange={onTaskTypeChange}
                    options={taskTypeOptions}
                    required={true}
                    error={taskDetailsError.type}
                    disabled={apiResponse.isLoading}
                />
                <Input
                    type='text'
                    placeholder='Title'
                    value={taskDetails.fields.title}
                    onChange={onTitleChange}
                    error={taskDetailsError.fields.title}
                    required={true}
                    disabled={apiResponse.isLoading}
                />
                <Input
                    type='text'
                    placeholder='Description'
                    value={taskDetails.fields.description}
                    onChange={onDescriptionChange}
                    disabled={apiResponse.isLoading}
                />
            </>
        </Modal>
    );
};

export default TaskModal;
