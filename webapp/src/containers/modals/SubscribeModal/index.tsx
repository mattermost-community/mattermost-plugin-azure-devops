import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import mm_constants from 'mattermost-redux/constants/general';

import Modal from 'components/modal';
import Form from 'components/form';
import EmptyState from 'components/emptyState';
import ResultPanel from 'components/resultPanel';

import {eventTypeBoards, eventTypeRepos, filterLabelValuePairAll} from 'pluginConstants/common';
import {boardEventTypeOptions, repoEventTypeOptions} from 'pluginConstants/form';
import pluginConstants from 'pluginConstants';

import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';

import {setServiceType, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowLinkModal} from 'reducers/linkModal';
import {getSubscribeModalState} from 'selectors';

import Utils from 'utils';

import ReposFilter from './filters/repos';

import './styles.scss';
import BoardsFilter from './filters/boards';

const SubscribeModal = () => {
    const {subscriptionModal} = pluginConstants.form;
    const [subscriptionModalFields, setSubscriptionModalFields] = useState<Record<SubscriptionModalFields, ModalFormFieldConfig>>(subscriptionModal);
    const [isFiltersError, setIsFiltersError] = useState<boolean>(false);

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
    const {visibility, project, organization, serviceType, projectID} = getSubscribeModalState(state);
    const {currentTeamId} = useSelector((reduxState: GlobalState) => reduxState.entities.teams);
    const {currentChannelId} = useSelector((reduxState: GlobalState) => reduxState.entities.channels);
    const dispatch = useDispatch();

    // State variables
    const [channelOptions, setChannelOptions] = useState<LabelValuePair[]>([]);
    const [showResultPanel, setShowResultPanel] = useState(false);

    // Function to hide the modal and reset all the states.
    const resetModalState = () => {
        dispatch(toggleShowSubscribeModal({isVisible: false, commandArgs: []}));
        resetFormFields();
        setChannelOptions([]);
        setShowResultPanel(false);
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
        const {isLoading, isSuccess, isError, data} = getApiState(
            pluginConstants.pluginApiServiceConfigs.getChannels.apiServiceName,
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

    useEffect(() => {
        if (formFields.serviceType === pluginConstants.common.boards) {
            setSubscriptionModalFields({...subscriptionModalFields, eventType: {...subscriptionModalFields.eventType, optionsList: boardEventTypeOptions}});
        } else if (formFields.serviceType === pluginConstants.common.repos) {
            setSubscriptionModalFields({...subscriptionModalFields, eventType: {...subscriptionModalFields.eventType, optionsList: repoEventTypeOptions}});
        }

        dispatch(setServiceType(formFields.serviceType ?? ''));
    }, [formFields.serviceType]);

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
        makeApiRequest(
            pluginConstants.pluginApiServiceConfigs.getChannels.apiServiceName,
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
        if (isChannelListSuccess && !showResultPanel) {
            setChannelOptions(channelList?.map((channel) => ({
                label: <span><i className={`icon ${channel.type === mm_constants.PRIVATE_CHANNEL ? 'icon-lock-outline' : 'icon-globe'} dropdown-option-icon`}/>{channel.display_name}</span>,
                value: channel.id,
            })));
        }

        // Pre-select the dropdown value in case of single option
        if (isOrganizationAndProjectListSuccess && !showResultPanel) {
            const autoSelectedValues: Pick<Record<FormFieldNames, string>, 'organization' | 'project' | 'channelID'> = {
                organization: organization ?? '',
                project: project ?? '',
                channelID: currentChannelId ?? '',
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

    const handleSetRepoFilter = (newValue: string, repoName?: string) =>
        setSpecificFieldValue({
            ...formFields,
            repository: newValue === filterLabelValuePairAll.value ? '' : newValue,
            repositoryName: repoName === filterLabelValuePairAll.label ? '' : repoName,
        });

    const handleSetTargetBranchFilter = (newValue: string) =>
        setSpecificFieldValue({
            ...formFields,
            targetBranch: newValue === filterLabelValuePairAll.value ? '' : newValue,
        });

    const handleSetPullRequestCreatedByFilter = (newValue: string, name?: string) =>
        setSpecificFieldValue({
            ...formFields,
            pullRequestCreatedBy: newValue,
            pullRequestCreatedByName: name === filterLabelValuePairAll.value ? '' : name,
        });

    const handleSetPullRequestReviewersContainsFilter = (newValue: string, name?: string) =>
        setSpecificFieldValue({
            ...formFields,
            pullRequestReviewersContains: newValue,
            pullRequestReviewersContainsName: name === filterLabelValuePairAll.value ? '' : name,
        });

    const handleSetPullRequestPushedByFilter = (newValue: string, name?: string) =>
        setSpecificFieldValue({
            ...formFields,
            pushedBy: newValue,
            pushedByName: name === filterLabelValuePairAll.value ? '' : name,
        });

    const handleSetPullRequestMergeResultFilter = (newValue: string, name?: string) =>
        setSpecificFieldValue({
            ...formFields,
            mergeResult: newValue,
            mergeResultName: name === filterLabelValuePairAll.value ? '' : name,
        });

    const handleSetPullRequestNotificationTypeFilter = (newValue: string, name?: string) =>
        setSpecificFieldValue({
            ...formFields,
            notificationType: newValue,
            notificationTypeName: name === filterLabelValuePairAll.value ? '' : name,
        });
    const handleSetAreaPathFilter = (newValue: string) =>
        setSpecificFieldValue({
            ...formFields,
            areaPath: newValue,
        });

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
                                            onChange={(newValue) => onChangeFormField(field as SubscriptionModalFields, newValue)}
                                            error={errorState[field as SubscriptionModalFields]}
                                            isDisabled={isLoading}
                                        />
                                    ))
                                }
                                {
                                    formFields.serviceType === pluginConstants.common.boards && formFields.eventType && Object.keys(eventTypeBoards).includes(formFields.eventType) && (
                                        <>
                                            <BoardsFilter
                                                organization={organization as string}
                                                projectId={projectID as string}
                                                eventType={formFields.eventType || ''}
                                                selectedAreaPath={formFields.areaPath || filterLabelValuePairAll.value}
                                                handleSelectAreaPath={handleSetAreaPathFilter}
                                                setIsFiltersError={setIsFiltersError}
                                            />
                                        </>
                                    )
                                }
                                {
                                    formFields.serviceType === pluginConstants.common.repos && formFields.eventType && Object.keys(eventTypeRepos).includes(formFields.eventType) && (
                                        <>
                                            <ReposFilter
                                                organization={organization as string}
                                                projectId={projectID as string}
                                                eventType={formFields.eventType || ''}
                                                selectedRepo={formFields.repository || filterLabelValuePairAll.value}
                                                handleSelectRepo={handleSetRepoFilter}
                                                selectedTargetBranch={formFields.targetBranch || filterLabelValuePairAll.value}
                                                handleSelectTargetBranch={handleSetTargetBranchFilter}
                                                selectedPullRequestCreatedBy={formFields.pullRequestCreatedBy || filterLabelValuePairAll.value}
                                                handleSelectPullRequestCreatedBy={handleSetPullRequestCreatedByFilter}
                                                selectedPullRequestReviewersContains={formFields.pullRequestReviewersContains || filterLabelValuePairAll.value}
                                                handlePullRequestReviewersContains={handleSetPullRequestReviewersContainsFilter}
                                                selectedPushedBy={formFields.pushedBy || filterLabelValuePairAll.value}
                                                handleSelectPushedBy={handleSetPullRequestPushedByFilter}
                                                selectedMergeResult={formFields.mergeResult || filterLabelValuePairAll.value}
                                                handleSelectMergeResult={handleSetPullRequestMergeResultFilter}
                                                selectedNotificationType={formFields.notificationType || filterLabelValuePairAll.value}
                                                handleSelectNotificationType={handleSetPullRequestNotificationTypeFilter}
                                                setIsFiltersError={setIsFiltersError}
                                            />
                                        </>
                                    )
                                }
                            </>
                        ) : (
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
