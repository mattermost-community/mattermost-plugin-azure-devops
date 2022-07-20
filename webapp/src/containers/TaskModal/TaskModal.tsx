import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideModal} from 'reducers/taskModal';

// TODO: fetch the organization and project options from API.
const organizationOptions = [
    {
        value: 'bs-test',
        label: 'bs-test',
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
    const [error, setError] = useState({
        taskOrganizationError: '',
        taskProjectError: '',
        taskTypeError: '',
        taskTitleError: '',
    });
    const [state, setState] = useState({
        taskOrganization: '',
        taskProject: '',
        taskType: '',
        taskTitle: '',
        taskDescription: '',
    });
    const [taskPayload, setTaskPayload] = useState<CreateTaskPayload | null>();
    const usePlugin = usePluginApi();
    const {visibility} = usePlugin.state['plugins-mattermost-plugin-azure-devops'].openTaskModalReducer;
    const dispatch = useDispatch();

    useEffect(() => {
        if (visibility === true) {
            // Pre-select the dropdown value in case of single option.
            if (organizationOptions.length === 1) {
                setState({...state, taskOrganization: organizationOptions[0].value});
            }
            if (projectOptions.length === 1) {
                setState({...state, taskProject: projectOptions[0].value});
            }
        }
    }, [visibility]);

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
        dispatch(hideModal());
    }, []);

    const onOrganizationChange = useCallback((value: string) => {
        setError({...error, taskOrganizationError: ''});
        setState({...state, taskOrganization: value});
    }, [error, state]);

    const onProjectChange = useCallback((value: string) => {
        setError({...error, taskProjectError: ''});
        setState({...state, taskProject: value});
    }, [error, state]);

    const onTaskTypeChange = useCallback((value: string) => {
        setError({...error, taskTypeError: ''});
        setState({...state, taskType: value});
    }, [error, state]);

    const onTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setError({...error, taskTitleError: ''});
        setState({...state, taskTitle: (e.target as HTMLInputElement).value});
    }, [error, state]);

    const onDescriptionChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setState({...state, taskDescription: (e.target as HTMLInputElement).value});
    }, [state]);

    const onConfirm = useCallback(() => {
        if (state.taskOrganization === '') {
            setError({...error, taskOrganizationError: 'Organization is required'});
            return;
        }
        if (state.taskProject === '') {
            setError({...error, taskProjectError: 'Project is required'});
            return;
        }
        if (state.taskType === '') {
            setError({...error, taskTypeError: 'Task type is required'});
            return;
        }
        if (state.taskTitle === '') {
            setError({...error, taskTitleError: 'Title is required'});
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
    }, [state, error]);

    useEffect(() => {
        if (taskPayload) {
            const {isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskPayload);
            if ((isSuccess && !isError) || (!isSuccess && isError)) {
                onHide();
            }
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
                loading={taskPayload ? usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskPayload).isLoading : false}
            >
                <>
                    <Dropdown
                        placeholder='Select organization'
                        value={state.taskOrganization}
                        onChange={(newValue) => onOrganizationChange(newValue)}
                        options={organizationOptions}
                        required={true}
                        error={error.taskOrganizationError}
                    />
                    <Dropdown
                        placeholder='Select project'
                        value={state.taskProject}
                        onChange={(newValue) => onProjectChange(newValue)}
                        options={projectOptions}
                        required={true}
                        error={error.taskProjectError}
                    />
                    <Dropdown
                        placeholder='Select work item'
                        value={state.taskType}
                        onChange={(newValue) => onTaskTypeChange(newValue)}
                        options={taskTypeOptions}
                        required={true}
                        error={error.taskTypeError}
                    />
                    <Input
                        type='text'
                        placeholder='Enter title'
                        value={state.taskTitle}
                        onChange={onTitleChange}
                        error={error.taskTitleError}
                        required={true}
                    />
                    <Input
                        type='text'
                        placeholder='Enter description'
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
