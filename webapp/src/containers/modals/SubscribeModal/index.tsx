import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import mm_constants from 'mattermost-redux/constants/general';

import {eventTypeBoards, eventTypePipelines, eventTypeRepos, filterLabelValuePairAll} from 'pluginConstants/common';
import {boardEventTypeOptions, pipelineEventTypeOptions, repoEventTypeOptions} from 'pluginConstants/form';
import pluginConstants from 'pluginConstants';

import Modal from 'components/modal';
import Form from 'components/form';
import EmptyState from 'components/emptyState';
import ResultPanel from 'components/resultPanel';

import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import usePluginApi from 'hooks/usePluginApi';
import useMattermostApi from 'hooks/useMattermostApi';
import useForm from 'hooks/useForm';

import {setServiceType, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowLinkModal} from 'reducers/linkModal';
import {getSubscribeModalState} from 'selectors';

import Utils from 'utils';

import ReposFilter from './filters/repos';
import BoardsFilter from './filters/boards';
import PipelinesFilter from './filters/pipelines';
import './styles.scss';

const SubscribeModal = () => {
    const {subscriptionModal} = pluginConstants.form;
    const [subscriptionModalFields, setSubscriptionModalFields] = useState<Record<SubscriptionModalFields, ModalFormFieldConfig>>(subscriptionModal);

    // Hooks
    const {
        formFields,
        errorState,
        onChangeFormField,
        setSpecificFieldValue,
        resetFormFields,
        isErrorInFormValidation,
    } = useForm(subscriptionModalFields);
    const {
        getApiState,
        makeApiRequest,
        makeApiRequestWithCompletionStatus,
        state,
    } = usePluginApi();
    const {makeMattermostApiRequest, getMattermostApiState} = useMattermostApi();
    const {visibility, project, organization, serviceType, projectID} = getSubscribeModalState(state);
    const {currentTeamId} = useSelector((reduxState: GlobalState) => reduxState.entities.teams);
    const {currentChannelId} = useSelector((reduxState: GlobalState) => reduxState.entities.channels);
    const dispatch = useDispatch();

    // State variables
    const [channelOptions, setChannelOptions] = useState<LabelValuePair[]>([]);
    const [showResultPanel, setShowResultPanel] = useState(false);
    const [isFiltersError, setIsFiltersError] = useState<boolean>(false);
    const [selectedProjectId, setSelectedProjectId] = useState<string>('');

    // Function to hide the modal and reset all the states.
    const resetModalState = () => {
        dispatch(toggleShowSubscribeModal({isVisible: false, commandArgs: []}));
        resetFormFields();
        setChannelOptions([]);
        setShowResultPanel(false);
        setIsFiltersError(false);
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

    // Get channel state
    const getChannelState = () => {
        const {isLoading, isSuccess, isError, data} = getMattermostApiState(
            pluginConstants.mattermostApiServiceConfigs.getChannels.apiServiceName,
            {teamId: currentTeamId},
        );
        return {isLoading, isSuccess, isError, data: data as ChannelList[]};
    };

    const {
        isSuccess: isOrganizationAndProjectListSuccess,
        isError: isOrganizationAndProjectListError,
        isLoading: isOrganizationAndProjectListLoading,
        organizationList,
        projectList,
    } = getOrganizationAndProjectState();

    const {data: channelList, isError: isChannelListError, isLoading: isChannelListLoading, isSuccess: isChannelListSuccess} = getChannelState();

    // Get option list for each types of dropdown fields
    const getDropDownOptions = (fieldName: SubscriptionModalFields) => {
        switch (fieldName) {
        case 'organization':
            return organizationList;
        case 'project':
            return projectList.filter(({metaData}) => metaData === formFields.organization);
        case 'serviceType':
            return subscriptionModalFields.serviceType.optionsList;
        case 'eventType':
            return subscriptionModalFields.eventType.optionsList;
        case 'channelID':
            return channelOptions;
        default:
            return [];
        }
    };

    const setSelectedDropdownOption = (field: SubscriptionModalFields, newValue: string, selectedOption?: Record<string, string>) => {
        onChangeFormField(field as SubscriptionModalFields, newValue);

        if (field === 'project' && selectedOption) {
            const selectedProject = selectedOption as ProjectListLabelValuePair;
            setSelectedProjectId(selectedProject.projectID ?? projectID);
        }
    };

    useEffect(() => {
        if (projectList.length === 1) {
            setSelectedProjectId(projectList[0].projectID);
        }
    }, [showResultPanel]);

    useEffect(() => {
        let optionsList: LabelValuePair[] = boardEventTypeOptions;

        if (formFields.serviceType === pluginConstants.common.repos) {
            optionsList = repoEventTypeOptions;
        } else if (formFields.serviceType === pluginConstants.common.pipelines) {
            optionsList = pipelineEventTypeOptions;
        }

        setSubscriptionModalFields({
            ...subscriptionModalFields,
            eventType: {...subscriptionModalFields.eventType, optionsList, isFieldDisabled: !formFields.project},
            serviceType: {...subscriptionModalFields.serviceType, isFieldDisabled: !formFields.project},
        });

        setSpecificFieldValue({
            ...formFields,
            eventType: optionsList[0].value,
        });

        dispatch(setServiceType(formFields.serviceType ?? ''));
    }, [formFields.serviceType, formFields.project]);

    // Opens link project modal
    const handleOpenLinkProjectModal = () => {
        dispatch(toggleShowLinkModal({isVisible: true, commandArgs: []}));
        resetModalState();
    };

    // Opens subscription modal
    const handleSubscriptionModal = () => {
        resetModalState();
        dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs: []}));
    };

    // Return different types of error messages occurred on API call
    const showApiErrorMessages = (isCreateSubscriptionError: boolean, error?: ApiErrorResponse) => {
        if (isChannelListError) {
            return pluginConstants.messages.error.errorFetchingChannelsList;
        }
        if (isOrganizationAndProjectListError) {
            return pluginConstants.messages.error.errorFetchingOrganizationAndProjectsList;
        }
        return Utils.getErrorMessage(isCreateSubscriptionError, 'SubscribeModal', error);
    };

    // Handles creating subscription on confirmation
    const onConfirm = () => {
        if (!isErrorInFormValidation()) {
            // Make POST api request to create subscription
            makeApiRequestWithCompletionStatus(
                pluginConstants.pluginApiServiceConfigs.createSubscription.apiServiceName,
                formFields as APIRequestPayload,
            );
        }
    };

    // Observe for the change in redux state after the API call to create a subscription and do the required actions
    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.createSubscription.apiServiceName,
        handleSuccess: () => {
            setShowResultPanel(true);
            dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs: [], isActionDone: true}));
        },
        payload: formFields as APIRequestPayload,
    });

    // Make API request to fetch channel list
    useEffect(() => {
        makeMattermostApiRequest(
            pluginConstants.mattermostApiServiceConfigs.getChannels.apiServiceName,
            {teamId: currentTeamId},
        );
    }, [visibility]);

    // Autoselect serviceType based on slash command
    useEffect(() => {
        setSpecificFieldValue({
            ...formFields,
            ...{serviceType},
        });
    }, [serviceType]);

    // Set organization, project and channel list values
    useEffect(() => {
        let isCurrentChannelIdPresentInChannelList = false; // Check if the current channel ID is the ID of a public or private channel and not the ID of a DM or group channel
        if (isChannelListSuccess && !showResultPanel) {
            const publicAndPrivateChannelList: LabelValuePair[] = [];

            if (channelList.length) {
                channelList.forEach((channel) => {
                    if (channel.type === mm_constants.PRIVATE_CHANNEL || channel.type === mm_constants.OPEN_CHANNEL) {
                        if (currentChannelId === channel.id) {
                            isCurrentChannelIdPresentInChannelList = true;
                        }

                        publicAndPrivateChannelList.push(({
                            label: <span><i className={`icon ${channel.type === mm_constants.PRIVATE_CHANNEL ? 'icon-lock-outline' : 'icon-globe'} azd-dropdown-option-icon`}/>{channel.display_name}</span>,
                            value: channel.id,
                        }));
                    }
                });
            }

            setChannelOptions(publicAndPrivateChannelList);
        }

        // Pre-select the dropdown value in case of single option
        if (isOrganizationAndProjectListSuccess && !showResultPanel) {
            const autoSelectedValues: Pick<Record<FormFieldNames, string>, 'organization' | 'project' | 'channelID'> = {
                organization: organization ?? '',
                project: project ?? '',
                channelID: isCurrentChannelIdPresentInChannelList && currentChannelId ? currentChannelId : '',
            };

            if (!organization && organizationList.length === 1) {
                autoSelectedValues.organization = organizationList[0].value;
            }
            if (!project && projectList.length === 1) {
                autoSelectedValues.project = projectList[0].value;
            }
            if (channelOptions.length === 1) {
                autoSelectedValues.channelID = channelOptions[0].value;
            }

            if (autoSelectedValues.organization || autoSelectedValues.project || autoSelectedValues.channelID) {
                setSpecificFieldValue({
                    ...formFields,
                    ...autoSelectedValues,
                });
            }
        }
    }, [
        isChannelListLoading,
        isOrganizationAndProjectListLoading,
        showResultPanel,
    ]);

    const handleSetSubscriptionFilter: HandleSetSubscriptionFilter = (
        filterID: FormFieldNames,
        filterIDNewValue: string,
        filterDisplayName?: FormFieldNames,
        filterDisplayNameNewValue?: string,
    ) => {
        let modifiedFields: Partial<Record<FormFieldNames, string>> = {
            [filterID]: filterIDNewValue === filterLabelValuePairAll.value ? '' : filterIDNewValue,
        };

        if (filterDisplayName && filterDisplayNameNewValue) {
            modifiedFields = {
                ...modifiedFields,
                [filterDisplayName]: filterDisplayNameNewValue === filterLabelValuePairAll.value ? '' : filterDisplayNameNewValue,
            };
        }

        if (filterID === 'repository') {
            modifiedFields = {
                ...modifiedFields,
                targetBranch: filterDisplayName === filterLabelValuePairAll.label ? '' : formFields.targetBranch,
            };
        }

        if (filterID === 'runPipeline') {
            modifiedFields = {
                ...modifiedFields,
                runStage: filterIDNewValue === filterLabelValuePairAll.value ? '' : formFields.runStage,
                runStageId: filterIDNewValue === filterLabelValuePairAll.value ? '' : formFields.runStageId,
            };
        }

        if (filterID === 'runStageStateId') {
            modifiedFields = {
                ...modifiedFields,
                runStageResultId: (filterIDNewValue !== filterLabelValuePairAll.value && filterIDNewValue !== 'Completed') ? '' : formFields.runStageResultId,
            };
        }

        if (filterID === 'runStateId') {
            modifiedFields = {
                ...modifiedFields,
                runResultId: (filterIDNewValue !== filterLabelValuePairAll.value && filterIDNewValue !== 'Completed') ? '' : formFields.runResultId,
            };
        }

        setSpecificFieldValue({
            ...formFields,
            ...modifiedFields,
        });
    };

    const {isLoading: isCreateSubscriptionLoading, isError, error} = getApiState(pluginConstants.pluginApiServiceConfigs.createSubscription.apiServiceName, formFields as APIRequestPayload);
    const isAnyProjectLinked = Boolean(organizationList.length && projectList.length);
    const isLoading = isChannelListLoading || isOrganizationAndProjectListLoading || isCreateSubscriptionLoading;

    return (
        <Modal
            show={visibility}
            title='Add New Subscription'
            onHide={resetModalState}
            onConfirm={isAnyProjectLinked ? onConfirm : null}
            confirmBtnText='Add New Subscription'
            confirmDisabled={isLoading}
            cancelDisabled={isLoading}
            loading={isLoading}
            showFooter={!showResultPanel}
            error={showApiErrorMessages(isError, error as ApiErrorResponse) || showApiErrorMessages(isFiltersError)}
        >
            <>
                {
                    !showResultPanel && (
                        isAnyProjectLinked ? (
                            <>
                                {
                                    Object.keys(subscriptionModalFields).map((field) => (
                                        <Form
                                            key={subscriptionModalFields[field as SubscriptionModalFields].label}
                                            fieldConfig={subscriptionModalFields[field as SubscriptionModalFields]}
                                            value={formFields[field as SubscriptionModalFields] ?? ''}
                                            optionsList={getDropDownOptions(field as SubscriptionModalFields)}
                                            onChange={(newValue, _, selectedOption) => setSelectedDropdownOption(field as SubscriptionModalFields, newValue, selectedOption)}
                                            error={errorState[field as SubscriptionModalFields]}
                                            isDisabled={isLoading}
                                        />
                                    ))
                                }
                                {
                                    formFields.serviceType === pluginConstants.common.boards && formFields.eventType && Object.keys(eventTypeBoards).includes(formFields.eventType) && (
                                        <BoardsFilter
                                            isModalOpen={visibility}
                                            organization={formFields.organization as string}
                                            projectId={selectedProjectId || projectID as string}
                                            eventType={formFields.eventType || ''}
                                            handleSetFilter={handleSetSubscriptionFilter}
                                            setIsFiltersError={setIsFiltersError}
                                            selectedAreaPath={formFields.areaPath || filterLabelValuePairAll.value}
                                        />
                                    )
                                }
                                {
                                    formFields.serviceType === pluginConstants.common.repos && formFields.eventType && Object.keys(eventTypeRepos).includes(formFields.eventType) && (
                                        <ReposFilter
                                            isModalOpen={visibility}
                                            organization={formFields.organization as string}
                                            projectId={selectedProjectId || projectID as string}
                                            eventType={formFields.eventType || ''}
                                            handleSetFilter={handleSetSubscriptionFilter}
                                            selectedRepo={formFields.repository || filterLabelValuePairAll.value}
                                            selectedTargetBranch={formFields.targetBranch || filterLabelValuePairAll.value}
                                            selectedPullRequestCreatedBy={formFields.pullRequestCreatedBy || filterLabelValuePairAll.value}
                                            selectedPullRequestReviewersContains={formFields.pullRequestReviewersContains || filterLabelValuePairAll.value}
                                            selectedPushedBy={formFields.pushedBy || filterLabelValuePairAll.value}
                                            selectedMergeResult={formFields.mergeResult || filterLabelValuePairAll.value}
                                            selectedNotificationType={formFields.notificationType || filterLabelValuePairAll.value}
                                            setIsFiltersError={setIsFiltersError}
                                        />
                                    )
                                }
                                {
                                    formFields.serviceType === pluginConstants.common.pipelines && formFields.eventType && Object.keys(eventTypePipelines).includes(formFields.eventType) && (
                                        <PipelinesFilter
                                            isModalOpen={visibility}
                                            organization={formFields.organization as string}
                                            projectId={selectedProjectId || projectID as string}
                                            eventType={formFields.eventType || ''}
                                            handleSetFilter={handleSetSubscriptionFilter}
                                            setIsFiltersError={setIsFiltersError}
                                            selectedBuildPipeline={formFields.buildPipeline || filterLabelValuePairAll.value}
                                            selectedBuildStatus={formFields.buildStatus || filterLabelValuePairAll.value}
                                            selectedReleasePipeline={formFields.releasePipeline || filterLabelValuePairAll.value}
                                            selectedStageName={formFields.stageName || filterLabelValuePairAll.value}
                                            selectedApprovalType={formFields.approvalType || filterLabelValuePairAll.value}
                                            selectedApprovalStatus={formFields.approvalStatus || filterLabelValuePairAll.value}
                                            selectedReleaseStatus={formFields.releaseStatus || filterLabelValuePairAll.value}
                                            selectedRunPipeline={formFields.runPipeline || filterLabelValuePairAll.value}
                                            selectedRunStage={formFields.runStage || filterLabelValuePairAll.value}
                                            selectedRunEnvironment={formFields.runEnvironment || filterLabelValuePairAll.value}
                                            selectedRunStageId={formFields.runStageId || filterLabelValuePairAll.value}
                                            selectedRunStageStateId={formFields.runStageStateId || filterLabelValuePairAll.value}
                                            selectedRunStageResultId={formFields.runStageResultId || filterLabelValuePairAll.value}
                                            selectedRunStateId={formFields.runStateId || filterLabelValuePairAll.value}
                                            selectedRunResultId={formFields.runResultId || filterLabelValuePairAll.value}
                                        />
                                    )
                                }
                            </>
                        ) :
                            !isLoading && (
                                <EmptyState
                                    title='No Project Linked'
                                    subTitle={{text: 'You can link a project by clicking the below button.'}}
                                    buttonText='Link new project'
                                    buttonAction={handleOpenLinkProjectModal}
                                />
                            )
                    )
                }
                {
                    showResultPanel && (
                        <ResultPanel
                            header='Subscription created successfully.'
                            primaryBtnText='Add New Subscription'
                            secondaryBtnText='Close'
                            onPrimaryBtnClick={handleSubscriptionModal}
                            onSecondaryBtnClick={resetModalState}
                        />)
                }
            </>
        </Modal>
    );
};

export default SubscribeModal;
