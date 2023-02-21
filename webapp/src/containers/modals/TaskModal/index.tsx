import React, {useCallback, useEffect, useMemo, useState} from 'react';
import {useDispatch} from 'react-redux';

import Modal from 'components/modal';
import EmptyState from 'components/emptyState';
import Form from 'components/form';
import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {toggleShowTaskModal} from 'reducers/taskModal';
import {toggleShowLinkModal} from 'reducers/linkModal';
import {getCreateTaskModalState} from 'selectors';

import pluginConstants from 'pluginConstants';
import {boardEventTypeOptions, subscriptionFiltersForBoards, subscriptionFiltersNameForBoards} from 'pluginConstants/form';

import Utils, {formLabelValuePairs} from 'utils';

const TaskModal = () => {
    const {createTaskModal: createTaskModalFields} = pluginConstants.form;

    // Hooks
    const {
        formFields,
        errorState,
        onChangeFormField,
        setSpecificFieldValue,
        resetFormFields,
        isErrorInFormValidation,
    } = useForm(createTaskModalFields);
    const {getApiState, makeApiRequestWithCompletionStatus, state} = usePluginApi();
    const dispatch = useDispatch();

    // State variables
    const {visibility, commandArgs} = getCreateTaskModalState(state);
    const [selectedProjectId, setSelectedProjectId] = useState<string>('');

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

    // Get organization and project state
    const getOrganizationAndProjectState = () => {
        const {isLoading, isSuccess, isError, data} = getApiState(
            pluginConstants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
        );

        return {
            isLoading,
            isError,
            isSuccess,
            organizationList: isSuccess ? Utils.getOrganizationList(data as ProjectDetails[]) : [],
            projectList: isSuccess ? Utils.getProjectList(data as ProjectDetails[]) : [],
        };
    };

    const getAreaPathValuesRequest = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization: formFields.organization as string,
        projectId: selectedProjectId,
        filters: subscriptionFiltersForBoards,
        eventType: boardEventTypeOptions[0].value, // We need an eventType field for this API. Since we do not have one in the create task modal so, we are hardcoding it.
    }), [formFields.organization, formFields.project, subscriptionFiltersForBoards, selectedProjectId]);

    useEffect(() => {
        if (visibility && getAreaPathValuesRequest.organization && getAreaPathValuesRequest.projectId) {
            makeApiRequestWithCompletionStatus(
                pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
                getAreaPathValuesRequest,
            );
        }
    }, [visibility, getAreaPathValuesRequest]);

    const {data: areaPathData, isLoading: isAreaPathLoading, isError: isAreaPathError, isSuccess: isAreaPathSuccess} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
        getAreaPathValuesRequest as APIRequestPayload,
    );

    const areaPathList = areaPathData as GetSubscriptionFiltersResponse || [];

    const getAreaPathOptions = useCallback(() => (isAreaPathSuccess ? ([...formLabelValuePairs('displayValue', 'value', areaPathList[subscriptionFiltersNameForBoards.areaPath], ['[Any]'])]) : []), [areaPathList]);

    const handleSetAreaPathField = (newValue: string) =>
        setSpecificFieldValue({
            ...formFields,
            areaPath: newValue,
        });

    const setSelectedDropdownOption = (field: CreateTaskModalFields, newValue: string, selectedOption?: Record<string, string>) => {
        onChangeFormField(field as CreateTaskModalFields, newValue);
        if (field === 'project' && selectedOption) {
            setSelectedProjectId((selectedOption as ProjectListLabelValuePair).projectID);
        }
    };

    const {
        isSuccess: isOrganizationAndProjectListSuccess,
        isError: isOrganizationAndProjectListError,
        isLoading: isOrganizationAndProjectListLoading,
        organizationList,
        projectList,
    } = getOrganizationAndProjectState();

    // Get option list for each types of dropdown fields
    const getDropDownOptions = (fieldName: CreateTaskModalFields) => {
        switch (fieldName) {
        case 'organization':
            return organizationList;
        case 'project':
            return projectList.filter((project) => project.metaData === formFields.organization);
        case 'type':
            return createTaskModalFields.type.optionsList;
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
                pluginConstants.pluginApiServiceConfigs.createTask.apiServiceName,
                getApiPayload(),
            );
        }
    };

    // Observe the change in redux state after the API call to create task and do the required actions
    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.createTask.apiServiceName,
        payload: getApiPayload(),
        handleSuccess: () => resetModalState(),
    });

    // Return different types of error messages occurred on API call
    const showApiErrorMessages = (isCreateSubscriptionError: boolean, error: ApiErrorResponse) => {
        if (isOrganizationAndProjectListError) {
            return pluginConstants.messages.error.errorFetchingOrganizationAndProjectsList;
        }
        return Utils.getErrorMessage(isCreateSubscriptionError, 'CreateTaskModal', error);
    };

    // Pre-select the dropdown value in case of single option
    useEffect(() => {
        if (isOrganizationAndProjectListSuccess) {
            const autoSelectedValues: Pick<Record<FormFieldNames, string>, 'organization' | 'project'> = {
                organization: '',
                project: '',
            };

            if (organizationList.length === 1) {
                autoSelectedValues.organization = organizationList[0].value;
            }
            if (projectList.length === 1) {
                autoSelectedValues.project = projectList[0].value;
                setSelectedProjectId(projectList[0].projectID);
            }

            if (autoSelectedValues.organization || autoSelectedValues.project) {
                setSpecificFieldValue({
                    ...formFields,
                    ...autoSelectedValues,
                });
            }
        }
    }, [isOrganizationAndProjectListLoading]);

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

    const {isLoading: isCreateTaskLoading, isError, error} = getApiState(pluginConstants.pluginApiServiceConfigs.createTask.apiServiceName, getApiPayload());
    const isAnyProjectLinked = Boolean(organizationList.length && projectList.length);
    const isLoading = isOrganizationAndProjectListLoading || isCreateTaskLoading || isAreaPathLoading;

    return (
        <Modal
            show={visibility}
            title='Create New Work Item'
            onHide={resetModalState}
            onConfirm={isAnyProjectLinked ? onConfirm : null}
            confirmBtnText='Create New Work Item'
            loading={isLoading}
            confirmDisabled={isLoading}
            error={showApiErrorMessages(isError, error as ApiErrorResponse)}
        >
            <>
                {
                    isAnyProjectLinked ? (
                        <>
                            {
                                Object.keys(createTaskModalFields).map((field) => (
                                    <Form
                                        key={createTaskModalFields[field as CreateTaskModalFields].label}
                                        fieldConfig={createTaskModalFields[field as CreateTaskModalFields]}
                                        value={formFields[field as CreateTaskModalFields] ?? ''}
                                        optionsList={getDropDownOptions(field as CreateTaskModalFields)}
                                        onChange={(newValue, _, selectedOption) => setSelectedDropdownOption(field as CreateTaskModalFields, newValue, selectedOption)}
                                        error={errorState[field as CreateTaskModalFields]}
                                        isDisabled={isLoading}
                                    />
                                ))
                            }
                            <Dropdown
                                placeholder='Area Path'
                                value={formFields.areaPath as string}
                                onChange={handleSetAreaPathField}
                                options={getAreaPathOptions()}
                                error={isAreaPathError}
                                loadingOptions={isAreaPathLoading}
                                disabled={!formFields.project || isLoading}
                            />
                        </>
                    ) : !isLoading && (
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
