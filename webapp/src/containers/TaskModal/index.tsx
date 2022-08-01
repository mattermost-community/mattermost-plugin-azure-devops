import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideTaskModal} from 'reducers/taskModal';
import {getOrganizationList, getProjectList} from 'utils';

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
    const [state, setState] = useState({
        taskOrganization: '',
        taskProject: '',
        taskType: '',
        taskTitle: '',
        taskDescription: '',
    });
    const [error, setError] = useState({
        taskOrganizationError: '',
        taskProjectError: '',
        taskTypeError: '',
        taskTitleError: '',
    });
    const [taskPayload, setTaskPayload] = useState<CreateTaskPayload | null>();
    const usePlugin = usePluginApi();
    const {visibility} = usePlugin.state['plugins-mattermost-plugin-azure-devops'].openTaskModalReducer;
    const [loading, setLoading] = useState(false);
    const [optionsLoading, setOptionsLoading] = useState(false);
    const [organizationOptions, setOrganizationOptions] = useState<DropdownOptionType[]>();
    const [projectOptions, setProjectOptions] = useState<DropdownOptionType[]>();
    const [APIError, setAPIError] = useState('');

    const dispatch = useDispatch();

    useEffect(() => {
        if (visibility === true) {
            usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.getProjects.apiServiceName);
        }
    }, [visibility]);

    useEffect(() => {
        // Pre-select the dropdown value in case of single option.
        if (organizationOptions?.length === 1) {
            setState((value) => ({...value, taskOrganization: organizationOptions[0].value}));
        }
        if (projectOptions?.length === 1) {
            setState((value) => ({...value, taskProject: projectOptions[0].value}));
        }
    }, [projectOptions, organizationOptions]);

    // Function to hide the modal and reset all the states.
    const onHide = useCallback(() => {
        setState({
            taskOrganization: '',
            taskProject: '',
            taskType: '',
            taskTitle: '',
            taskDescription: '',
        });
        setError({
            taskOrganizationError: '',
            taskProjectError: '',
            taskTypeError: '',
            taskTitleError: '',
        });
        setTaskPayload(null);
        dispatch(hideTaskModal());
    }, []);

    const onOrganizationChange = useCallback((value: string) => {
        setError({...error, taskOrganizationError: ''});
        setState({...state, taskOrganization: value});
    }, [state, error]);

    const onProjectChange = useCallback((value: string) => {
        setError({...error, taskProjectError: ''});
        setState({...state, taskProject: value});
    }, [state, error]);

    const onTaskTypeChange = useCallback((value: string) => {
        setError({...error, taskTypeError: ''});
        setState({...state, taskType: value});
    }, [state, error]);

    const onTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setError({...error, taskTitleError: ''});
        setState({...state, taskTitle: (e.target as HTMLInputElement).value});
    }, [state, error]);

    const onDescriptionChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setState({...state, taskDescription: (e.target as HTMLInputElement).value});
    }, [state]);

    const onConfirm = useCallback(() => {
        if (state.taskOrganization === '') {
            setError((value) => ({...value, taskOrganizationError: 'Organization is required'}));
        }
        if (state.taskProject === '') {
            setError((value) => ({...value, taskProjectError: 'Project is required'}));
        }
        if (state.taskType === '') {
            setError((value) => ({...value, taskTypeError: 'Work item type is required'}));
        }
        if (state.taskTitle === '') {
            setError((value) => ({...value, taskTitleError: 'Title is required'}));
        }

        if (!state.taskOrganization || !state.taskProject || !state.taskTitle || !state.taskType) {
            return;
        }

        // Create payload to send in the POST request.
        const payload = {
            organization: state.taskOrganization,
            project: state.taskProject,
            type: state.taskType,
            fields: {
                title: state.taskTitle,
                description: state.taskDescription,
            },
        };

        // TODO: save the payload in a state variable to use it while reading the state
        // we can see later if there exists a better way for this
        setTaskPayload(payload);

        // Make POST api request
        usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.createTask.apiServiceName, payload);
    }, [state]);

    useEffect(() => {
        if (taskPayload) {
            const {isLoading, isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskPayload);
            setLoading(isLoading);
            if (isSuccess) {
                onHide();
            }
            if (isError) {
                setAPIError('Failed to create the task.');
            }
            return;
        }
        const {data, isLoading, isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.getProjects.apiServiceName);
        setOptionsLoading(isLoading);
        if (isSuccess && data) {
            setProjectOptions(getProjectList(data));
            setOrganizationOptions(getOrganizationList(data));
        }
        if (isError) {
            setAPIError('Failed to load the options.');
        }
    }, [usePlugin.state]);

    if (visibility) {
        return (
            <Modal
                show={visibility}
                title='Create Task'
                onHide={onHide}
                onConfirm={onConfirm}
                confirmBtnText='Create task'
                loading={loading}
                confirmDisabled={loading}
                cancelDisabled={loading}
                error={APIError}
            >
                <>
                    <Dropdown
                        placeholder='Organization name'
                        value={state.taskOrganization}
                        onChange={(newValue) => onOrganizationChange(newValue)}
                        options={organizationOptions ?? []}
                        loadingOptions={optionsLoading}
                        required={true}
                        error={error.taskOrganizationError}
                    />
                    <Dropdown
                        placeholder='Project name'
                        value={state.taskProject}
                        onChange={(newValue) => onProjectChange(newValue)}
                        options={projectOptions ?? []}
                        loadingOptions={optionsLoading}
                        required={true}
                        error={error.taskProjectError}
                    />
                    <Dropdown
                        placeholder='Work item type'
                        value={state.taskType}
                        onChange={(newValue) => onTaskTypeChange(newValue)}
                        options={taskTypeOptions}
                        required={true}
                        error={error.taskTypeError}
                    />
                    <Input
                        type='text'
                        placeholder='Title'
                        value={state.taskTitle}
                        onChange={onTitleChange}
                        error={error.taskTitleError}
                        required={true}
                    />
                    <Input
                        type='text'
                        placeholder='Description'
                        value={state.taskDescription}
                        onChange={onDescriptionChange}
                    />
                </>
            </Modal>
        );
    }
    return null;
};

export default TaskModal;
