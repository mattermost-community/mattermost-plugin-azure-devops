import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideTaskModal} from 'reducers/taskModal';

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
    const [state, setState] = useState({
        taskOrganization: '',
        taskProject: '',
        taskType: '',
        taskTitle: '',
        taskDescription: '',
    });
    const [taskOrganizationError, setTaskOrganizationError] = useState('');
    const [taskProjectError, setTaskProjectError] = useState('');
    const [taskTypeError, setTaskTypeError] = useState('');
    const [taskTitleError, setTaskTitleError] = useState('');
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
        setTaskOrganizationError('');
        setTaskProjectError('');
        setTaskTitleError('');
        setTaskTypeError('');
        setTaskPayload(null);
        dispatch(hideTaskModal());
    }, []);

    const onOrganizationChange = useCallback((value: string) => {
        setTaskOrganizationError('');
        setState({...state, taskOrganization: value});
    }, [state]);

    const onProjectChange = useCallback((value: string) => {
        setTaskProjectError('');
        setState({...state, taskProject: value});
    }, [state]);

    const onTaskTypeChange = useCallback((value: string) => {
        setTaskTypeError('');
        setState({...state, taskType: value});
    }, [state]);

    const onTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setTaskTitleError('');
        setState({...state, taskTitle: (e.target as HTMLInputElement).value});
    }, [state]);

    const onDescriptionChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setState({...state, taskDescription: (e.target as HTMLInputElement).value});
    }, [state]);

    const onConfirm = useCallback(() => {
        if (state.taskOrganization === '') {
            setTaskOrganizationError('Organization is required');
        }
        if (state.taskProject === '') {
            setTaskProjectError('Project is required');
        }
        if (state.taskType === '') {
            setTaskTypeError('Work item type is required');
        }
        if (state.taskTitle === '') {
            setTaskTitleError('Title is required');
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
                        placeholder='Organization name'
                        value={state.taskOrganization}
                        onChange={(newValue) => onOrganizationChange(newValue)}
                        options={organizationOptions}
                        required={true}
                        error={taskOrganizationError}
                    />
                    <Dropdown
                        placeholder='Project name'
                        value={state.taskProject}
                        onChange={(newValue) => onProjectChange(newValue)}
                        options={projectOptions}
                        required={true}
                        error={taskProjectError}
                    />
                    <Dropdown
                        placeholder='Work item type'
                        value={state.taskType}
                        onChange={(newValue) => onTaskTypeChange(newValue)}
                        options={taskTypeOptions}
                        required={true}
                        error={taskTypeError}
                    />
                    <Input
                        type='text'
                        placeholder='Title'
                        value={state.taskTitle}
                        onChange={onTitleChange}
                        error={taskTitleError}
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
