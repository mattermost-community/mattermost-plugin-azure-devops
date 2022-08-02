import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Dropdown from 'components/dropdown';
import Input from 'components/inputField';
import Modal from 'components/modal';

import Constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';
import {hideTaskModal} from 'reducers/taskModal';
import {getOrganizationList, getProjectList} from 'utils';
import {getTaskModalState, getUserConnectionState} from 'selectors';
import LinearLoader from 'components/loader/linear';

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
    const [optionsLoading, setOptionsLoading] = useState(false);

    // Hooks
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    useEffect(() => {
        if (getTaskModalState(usePlugin.state).visibility && !usePlugin.getUserAccountConnectionState().isLoading &&
            usePlugin.getUserAccountConnectionState().isSuccess &&
            usePlugin.getUserAccountConnectionState().data?.MattermostUserID) {
            // Make API request to fetch all linked projects.
            usePlugin.makeApiRequest(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        }
    }, [usePlugin.state]);

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
        const newErrorState = {
            organization: '',
            project: '',
            type: '',
            title: '',
        };

        if (taskDetails.organization === '') {
            newErrorState.organization = 'Project is required';
        }
        if (taskDetails.project === '') {
            newErrorState.project = 'Project is required';
        }
        if (taskDetails.type === '') {
            newErrorState.type = 'Work item type is required';
        }
        if (taskDetails.fields.title === '') {
            newErrorState.title = 'Title is required';
        }

        if (newErrorState.organization || newErrorState.project || newErrorState.title || newErrorState.type) {
            setErrorState(newErrorState);
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

    return (
        <Modal
            show={getTaskModalState(usePlugin.state).visibility && !optionsLoading && !usePlugin.getApiState(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isError}
            title='Create Task'
            onHide={resetModalState}
            onConfirm={onConfirm}
            showFooter={
                !usePlugin.getUserAccountConnectionState().isLoading &&
                usePlugin.getUserAccountConnectionState().isSuccess &&
                usePlugin.getUserAccountConnectionState().data?.MattermostUserID
            }
            confirmBtnText='Create task'
            loading={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
            confirmDisabled={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
            cancelDisabled={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isLoading}
            error={usePlugin.getApiState(Constants.pluginApiServiceConfigs.createTask.apiServiceName, taskDetails).isError ? 'Failed to create the task' : ''}
        >
            <>
                {
                    usePlugin.getUserAccountConnectionState().isLoading && (<LinearLoader/>)
                }
                {
                    !usePlugin.getUserAccountConnectionState().isLoading &&
                    usePlugin.getUserAccountConnectionState().isError &&
                    (<div className='not-linked'>{'You do not have any Azure Devops account connected. Kindly link the account first'}</div>)
                }
                {
                    !usePlugin.getUserAccountConnectionState().isLoading &&
                    usePlugin.getUserAccountConnectionState().isSuccess &&
                    usePlugin.getUserAccountConnectionState().data?.MattermostUserID &&
                    usePlugin.getApiState(Constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).data?.length <= 0 &&
                    (<div className='not-linked'>{'You do not have any linked project. Kindly link a project first'}</div>)
                }
                {
                    !usePlugin.getUserAccountConnectionState().isLoading &&
                    usePlugin.getUserAccountConnectionState().isSuccess &&
                    usePlugin.getUserAccountConnectionState().data?.MattermostUserID && (
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
                    )
                }
            </>
        </Modal>
    );
};

export default TaskModal;
