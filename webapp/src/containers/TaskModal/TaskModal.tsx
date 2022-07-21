import React, {useCallback, useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Modal from 'components/modal';
import Input from 'components/inputField';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideModal} from 'reducers/taskModal';

const TaskModal = () => {
    const usePlugin = usePluginApi();
    const {visibility, title, description} = usePlugin.state['plugins-mattermost-plugin-azure-devops'].openTaskModalReducer;
    const [visible, setVisible] = useState(false);
    const [error, setError] = useState('');
    const [taskTitle, setTaskTitle] = useState('');
    const [taskDescription, setTaskDescription] = useState('');
    const dispatch = useDispatch();

    useEffect(() => {
        if (visibility === true) {
            setVisible(true);
            setTaskTitle(title || '');
            setTaskDescription(description || '');
        }
    }, [visibility]);

    const onHide = useCallback(() => {
        setVisible(false);
        setTaskTitle('');
        setTaskDescription('');
        setError('');
        dispatch(hideModal());
    }, []);

    const onTitleChange = useCallback((e: React.ChangeEvent<Element>) => {
        setError('');
        setTaskTitle((e.target as HTMLInputElement).value);
    }, []);

    const onDescriptionChange = useCallback((e: React.ChangeEvent<Element>) => {
        setTaskDescription((e.target as HTMLInputElement).value);
    }, []);

    const [taskPayload, setTaskPayload] = useState<CreateTaskPayload>()

    const onConfirm = () => {
        if (taskTitle === '') {
            setError('Title is required');
            return;
        }

        const payload = {
            organization: 'brightscout-test',
            project: 'azure-test',
            type: 'task',
            fields: {
                title: taskTitle,
                description: taskDescription
            }
        }

        // TODO: save the payload in a state variable to use it while reading the state
        // we can see later if there exists a better way for this
        setTaskPayload(payload)

        // Make POST api request
        usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.createTask.apiServiceName, payload);
    };

    // TODO: example, make a GET request, remove later if not required
    useEffect(() => {
        usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.testGet.apiServiceName)
    }, [])

    // TODO: example, reading GET & POST state changes, remove later if not required
    useEffect(() => {
        if (taskPayload) {
            const {data, isLoading, isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskPayload);
            console.log("Create-task state from component", data, isLoading, isSuccess, isError, usePlugin.state, taskPayload);
        }

        const {data, isLoading, isSuccess, isError} = usePlugin.getApiState(Constants.pluginApiServiceConfigs.testGet.apiServiceName);
        console.log("get-task state from component", data, isLoading, isSuccess, isError, usePlugin.state, taskPayload);

    }, [usePlugin.state, taskPayload])

    const input = (
        <>
            <Input
                type='text'
                placeholder='Enter title'
                label='Title'
                required={true}
                value={taskTitle}
                onChange={onTitleChange}
                error={error || ''}
            />
            <Input
                type='text'
                placeholder='Enter description'
                label='Description'
                value={taskDescription}
                onChange={onDescriptionChange}
            />
        </>
    );

    if (visible) {
        return (
            <>
                <Modal
                    show={visible}
                    title='Create Task'
                    onHide={onHide}
                    onConfirm={onConfirm}
                    confirmBtnText='Create task'
                    children={input}
                />
            </>
        );
    }
    return null;
};

export default TaskModal;
