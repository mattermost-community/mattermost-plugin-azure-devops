import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import Modal from 'components/modal';
import CircularLoader from 'components/loader/circular';
import EmptyState from 'components/emptyState';
import Form from 'components/form';

import plugin_constants from 'plugin_constants';

import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {toggleShowTaskModal} from 'reducers/taskModal';
import {toggleShowLinkModal} from 'reducers/linkModal';
import {getCreateTaskModalState} from 'selectors';

import Utils from 'utils';

const TaskModal = () => {
    const {createTaskModal} = plugin_constants.form;

    // Hooks
    const {
        formFields,
        errorState,
        onChangeFormField,
        setSpecificFieldValue,
        resetFormFields,
        isErrorInFormValidation,
    } = useForm(createTaskModal);
    const {getApiState, makeApiRequestWithCompletionStatus, state} = usePluginApi();
    const dispatch = useDispatch();

    // State variables
    const [organizationOptions, setOrganizationOptions] = useState<LabelValuePair[]>([]);
    const [projectOptions, setProjectOptions] = useState<LabelValuePair[]>([]);
    const {visibility, commandArgs} = getCreateTaskModalState(state);

    // Function to hide the modal and reset all the states
    const resetModalState = () => {
        dispatch(toggleShowTaskModal({isVisible: false, commandArgs: []}));
        resetFormFields();
    };

    // Opens link project modal
    const handleOpenLinkProjectModal = () => {
        dispatch(toggleShowLinkModal({isVisible: true, commandArgs: []}));
        resetModalState();
    };

    // Get option list for each types of dropdown fields
    const getDropDownOptions = (fieldName: CreateTaskModalFields) => {
        switch (fieldName) {
        case 'organization':
            return organizationOptions;
        case 'project':
            return projectOptions;
        case 'type':
            return createTaskModal.type.optionsList;
        default:
            return [];
        }
    };

    // Form payload to send in API request
    const getApiPayload = (): CreateTaskPayload => {
        const payload: CreateTaskPayload = {
            organization: formFields.organization ?? '',
            project: formFields.project ?? '',
            type: formFields.type ?? '',
            fields: {
                title: formFields.title ?? '',
                description: formFields.description ?? '',
                areaPath: formFields.areaPath ?? '',
            },
            timestamp: formFields.timestamp ?? '',
        };

        return payload;
    };

    // Handles creating a new task on confirmation
    const onConfirm = () => {
        if (!isErrorInFormValidation()) {
            // Make POST api request
            makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.createTask.apiServiceName,
                getApiPayload(),
            );
        }
    };

    // Observe the change in redux state after the API call to create task and do the required actions
    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.createTask.apiServiceName,
        payload: getApiPayload(),
        handleSuccess: () => resetModalState(),
    });

    // Get organization and project state
    const getOrganizationAndProjectState = () => {
        const {isLoading, isSuccess, isError, data} = getApiState(
            plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
        );

        return {
            isLoading,
            isError,
            isSuccess,
            organizationList: isSuccess ? Utils.getOrganizationList(data as ProjectDetails[]) : [],
            projectList: isSuccess ? Utils.getProjectList(data as ProjectDetails[]) : [],
        };
    };

    // Return different types of error messages occurred on API call
    const showApiErrorMessages = (isCreateSubscriptionError: boolean, error: ApiErrorResponse) => {
        if (getOrganizationAndProjectState().isError) {
            return plugin_constants.messages.error.errorFetchingOrganizationAndProjectsList;
        }
        return Utils.getErrorMessage(isCreateSubscriptionError, 'CreateTaskModal', error);
    };

    // Pre-select the dropdown value in case of single option
    useEffect(() => {
        const autoSelectedValues: Pick<Record<FormFieldNames, string>, 'organization' | 'project'> = {
            organization: '',
            project: '',
        };

        if (organizationOptions.length === 1) {
            autoSelectedValues.organization = organizationOptions[0].value;
        }
        if (projectOptions.length === 1) {
            autoSelectedValues.project = projectOptions[0].value;
        }

        if (autoSelectedValues.organization || autoSelectedValues.project) {
            setSpecificFieldValue({
                ...formFields,
                ...autoSelectedValues,
            });
        }
    }, [projectOptions, organizationOptions]);

    // Set organization and project list values
    useEffect(() => {
        if (getOrganizationAndProjectState().isSuccess) {
            setOrganizationOptions(getOrganizationAndProjectState().organizationList);
            setProjectOptions(getOrganizationAndProjectState().projectList);
        }
    }, [
        getOrganizationAndProjectState().isLoading,
    ]);

    // Set modal field values
    useEffect(() => {
        if (visibility) {
            setSpecificFieldValue({
                ...formFields,
                title: commandArgs.title,
                description: commandArgs.description,
            });
        }
    }, [visibility]);

    const {isLoading, isError, error} = getApiState(plugin_constants.pluginApiServiceConfigs.createTask.apiServiceName, getApiPayload());
    const isAnyProjectLinked = Boolean(getOrganizationAndProjectState().organizationList.length && getOrganizationAndProjectState().projectList.length);

    return (
        <Modal
            show={visibility}
            title='Create New Task'
            onHide={resetModalState}
            onConfirm={isAnyProjectLinked ? onConfirm : null}
            confirmBtnText='Create new task'
            loading={isLoading}
            confirmDisabled={isLoading}
            error={showApiErrorMessages(isError, error as ApiErrorResponse)}
        >
            <>
                {
                    getOrganizationAndProjectState().isLoading && <CircularLoader/>
                }
                {
                    isAnyProjectLinked ? (
                        Object.keys(createTaskModal).map((field) => (
                            <Form
                                key={createTaskModal[field as CreateTaskModalFields].label}
                                fieldConfig={createTaskModal[field as CreateTaskModalFields]}
                                value={formFields[field as CreateTaskModalFields] ?? ''}
                                optionsList={getDropDownOptions(field as CreateTaskModalFields)}
                                onChange={(newValue) => onChangeFormField(field as CreateTaskModalFields, newValue)}
                                error={errorState[field as CreateTaskModalFields]}
                                isDisabled={isLoading}
                            />
                        ))
                    ) : (
                        <EmptyState
                            title='No Project Linked'
                            subTitle={{text: 'You can link a project by clicking the below button.'}}
                            buttonText='Link New Project'
                            buttonAction={handleOpenLinkProjectModal}
                        />
                    )
                }
            </>
        </Modal>
    );
};

export default TaskModal;
