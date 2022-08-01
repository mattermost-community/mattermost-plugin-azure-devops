import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideTaskModal} from 'reducers/taskModal';
import {getOrganizationList, getProjectList} from 'utils';
import {getTaskModalState} from 'selectors';

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
    const [errorState, setErrorState] = useState({
        organization: '',
        project: '',
        type: '',
        title: '',
    });
    const [organizationOptions, setOrganizationOptions] = useState<DropdownOptionType[]>();
    const [projectOptions, setProjectOptions] = useState<DropdownOptionType[]>();

    // const [loading, setLoading] = useState(false);
    const [optionsLoading, setOptionsLoading] = useState(false);

    // Hooks
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    const {visibility} = getTaskModalState(usePlugin.state);

    useEffect(() => {
        if (visibility === true) {
            // Make API request to fetch all linked projects.
            usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        }
    }, [visibility]);

    useEffect(() => {
        // Pre-select the dropdown value in case of single option.
        if (organizationOptions?.length === 1) {
            setTaskDetails((value) => ({...value, organization: organizationOptions[0].value}));
        }
        if (projectOptions?.length === 1) {
            setTaskDetails((value) => ({...value, project: projectOptions[0].value}));
        }
    }, [projectOptions, organizationOptions]);

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
        setErrorState({
            organization: '',
            project: '',
            type: '',
            title: '',
        });
        setOrganizationOptions([]);
        setProjectOptions([]);
        setOptionsLoading(false);

        // setLoading(false);
        dispatch(hideTaskModal());
    };

    // Set organization name.
    const onOrganizationChange = (value: string) => {
        setErrorState({...errorState, organization: ''});
        setTaskDetails({...taskDetails, organization: value});
    };

    // Set project name.
    const onProjectChange = (value: string) => {
        setErrorState({...errorState, project: ''});
        setTaskDetails({...taskDetails, project: value});
    };

    // Set task type.
    const onTaskTypeChange = (value: string) => {
        setErrorState({...errorState, type: ''});
        setTaskDetails({...taskDetails, type: value});
    };

    // Set task title.
    const onTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setErrorState({...errorState, title: ''});
        setTaskDetails({...taskDetails, fields: {...taskDetails.fields, title: (e.target as HTMLInputElement).value}});
    };

    // Set task description.
    const onDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTaskDetails({...taskDetails, fields: {...taskDetails.fields, description: (e.target as HTMLInputElement).value}});
    };

    const onConfirm = () => {
        if (taskDetails.organization === '') {
            errorState.organization = 'Project is required';
        }
        if (taskDetails.project === '') {
            errorState.project = 'Project is required';
        }
        if (taskDetails.type === '') {
            errorState.type = 'Work item type is required';
        }
        if (taskDetails.fields.title === '') {
            errorState.title = 'Title is required';
        }

        if (errorState.organization || errorState.project || errorState.title || errorState.type) {
            return;
        }

        // Make POST api request
        createTask(taskDetails);
    };

    // Make POST api request to link a project
    const createTask = async (payload: CreateTaskPayload) => {
        const createTaskResponse = await usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.createTask.apiServiceName, payload);
        if (createTaskResponse) {
            resetModalState();
        }
    };

    useEffect(() => {
        const {data, isLoading, isSuccess} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        setOptionsLoading(isLoading);
        if (isSuccess && data) {
            setProjectOptions(getProjectList(data));
            setOrganizationOptions(getOrganizationList(data));
        }
    }, [usePlugin.state]);

    if (visibility) {
        return (
            <Modal
                show={visibility && !optionsLoading && !usePlugin.getApiState(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isError}
                title='Create Task'
                onHide={resetModalState}
                onConfirm={onConfirm}
                confirmBtnText='Create task'
                loading={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
                confirmDisabled={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
                cancelDisabled={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
                error={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isError ? 'Failed to create the task' : ''}
            >
                <>
                    <Dropdown
                        placeholder='Organization name'
                        value={taskDetails.organization}
                        onChange={(newValue) => onOrganizationChange(newValue)}
                        options={organizationOptions ?? []}
                        loadingOptions={optionsLoading}
                        required={true}
                        error={errorState.organization}
                    />
                    <Dropdown
                        placeholder='Project name'
                        value={taskDetails.project}
                        onChange={(newValue) => onProjectChange(newValue)}
                        options={projectOptions ?? []}
                        loadingOptions={optionsLoading}
                        required={true}
                        error={errorState.project}
                    />
                    <Dropdown
                        placeholder='Work item type'
                        value={taskDetails.type}
                        onChange={(newValue) => onTaskTypeChange(newValue)}
                        options={taskTypeOptions}
                        required={true}
                        error={errorState.type}
                    />
                    <Input
                        type='text'
                        placeholder='Title'
                        value={taskDetails.fields.title}
                        onChange={onTitleChange}
                        error={errorState.title}
                        required={true}
                    />
                    <Input
                        type='text'
                        placeholder='Description'
                        value={taskDetails.fields.description}
                        onChange={onDescriptionChange}
                    />
                </>
            </Modal>
        );
    }
    return null;
};

export default TaskModal;
