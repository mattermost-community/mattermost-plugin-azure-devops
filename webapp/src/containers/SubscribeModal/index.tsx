import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';

import Modal from 'components/modal';
import CircularLoader from 'components/loader/circular';
import Form from 'components/form';
import EmptyState from 'components/emptyState';

import plugin_constants from 'plugin_constants';

import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';

import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowLinkModal} from 'reducers/linkModal';
import {getSubscribeModalState} from 'selectors';

import Utils from 'utils';

import './styles.scss';

const SubscribeModal = () => {
    // Hooks
    const {
        formFields,
        errorState,
        onChangeOfFormField,
        setSpecificFieldValue,
        resetFormFields,
        isErrorInFormValidation,
    } = useForm(plugin_constants.form.subscriptionModal);
    const {
        getApiState,
        makeApiRequest,
        makeApiRequestWithCompletionStatus,
        state,
    } = usePluginApi();
    const {visibility} = getSubscribeModalState(state);
    const {entities} = useSelector((reduxState: GlobalState) => reduxState);
    const dispatch = useDispatch();

    // State variables
    const [channelOptions, setChannelOptions] = useState<LabelValuePair[]>([]);
    const [organizationOptions, setOrganizationOptions] = useState<LabelValuePair[]>([]);
    const [projectOptions, setProjectOptions] = useState<LabelValuePair[]>([]);

    // Function to hide the modal and reset all the states.
    const resetModalState = (isActionDone?: boolean) => {
        setChannelOptions([]);
        setOrganizationOptions([]);
        setProjectOptions([]);
        resetFormFields();
        dispatch(toggleShowSubscribeModal({isVisible: false, commandArgs: [], isActionDone}));
    };

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

    // Get channel state
    const getChannelState = () => {
        const {isLoading, isSuccess, isError, data} = getApiState(
            plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName,
            {teamId: entities.teams.currentTeamId},
        );
        return {isLoading, isSuccess, isError, data: data as ChannelList[]};
    };

    // Get option list for each types of dropdown fields
    const getDropDownOptions = (fieldName: SubscriptionModalFields) => {
        switch (fieldName) {
        case 'organization':
            return organizationOptions;
        case 'project':
            return projectOptions;
        case 'eventType':
            return plugin_constants.form.subscriptionModal.eventType.optionsList;
        case 'channelID':
            return channelOptions;
        default:
            return [];
        }
    };

    // Opens link project modal
    const handleOpenLinkProjectModal = () => {
        dispatch(toggleShowLinkModal({isVisible: true, commandArgs: []}));
        resetModalState();
    };

    // Return different types of error messages occurred on API call
    const showApiErrorMessages = (isCreateSubscriptionError: boolean, error: ApiErrorResponse) => {
        if (getChannelState().isError) {
            return plugin_constants.messages.error.errorFetchingChannelsList;
        }
        if (getOrganizationAndProjectState().isError) {
            return plugin_constants.messages.error.errorFetchingOrganizationAndProjectsList;
        }
        return Utils.getErrorMessage(isCreateSubscriptionError, 'SubscribeModal', error);
    };

    // Handles on confirming create subscription
    const onConfirm = () => {
        if (!isErrorInFormValidation()) {
            // Make POST api request to create subscription
            makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName,
                formFields as APIRequestPayload,
            );
        }
    };

    // Observe for the change in redux state after API call and do the required actions
    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName,
        handleSuccess: () => resetModalState(true),
        payload: formFields as APIRequestPayload,
    });

    // Make API request to fetch channel list
    useEffect(() => {
        makeApiRequest(
            plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName,
            {teamId: entities.teams.currentTeamId},
        );
    }, [visibility]);

    // Pre-select the dropdown value in case of single option
    useEffect(() => {
        const autoSelectedValues: Pick<Record<FormFields, string>, 'organization' | 'project' | 'channelID'> = {
            organization: '',
            project: '',
            channelID: '',
        };

        if (organizationOptions.length === 1) {
            autoSelectedValues.organization = organizationOptions[0].value;
        }
        if (projectOptions.length === 1) {
            autoSelectedValues.project = projectOptions[0].value;
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
    }, [projectOptions, organizationOptions, channelOptions]);

    // Set organization, project and channel list values
    useEffect(() => {
        if (getChannelState().isSuccess) {
            setChannelOptions(getChannelState().data?.map((channel) => ({
                label: <span><i className='fa fa-globe dropdown-option-icon'/>{channel.display_name}</span>,
                value: channel.id,
            })));
        }

        if (getOrganizationAndProjectState().isSuccess) {
            setOrganizationOptions(getOrganizationAndProjectState().organizationList);
            setProjectOptions(getOrganizationAndProjectState().projectList);
        }
    }, [
        getChannelState().isLoading,
        getOrganizationAndProjectState().isLoading,
    ]);

    const {isLoading, isError, error} = getApiState(plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName, formFields as APIRequestPayload);
    const isAnyProjectLinked = Boolean(getOrganizationAndProjectState().organizationList.length && getOrganizationAndProjectState().projectList.length);

    return (
        <Modal
            show={visibility}
            title='Create subscription'
            onHide={resetModalState}
            onConfirm={isAnyProjectLinked ? onConfirm : null}
            confirmBtnText='Create subscription'
            confirmDisabled={isLoading}
            cancelDisabled={isLoading}
            loading={isLoading}
            error={showApiErrorMessages(isError, error as ApiErrorResponse)}
        >
            <>
                {
                    (getChannelState().isLoading || getOrganizationAndProjectState().isLoading) && <CircularLoader/>
                }
                {
                    !isAnyProjectLinked && (
                        <EmptyState
                            title='No Project Linked'
                            subTitle={{text: 'Link a project by clicking the button below'}}
                            buttonText='Link new project'
                            buttonAction={handleOpenLinkProjectModal}
                        />
                    )
                }
                {
                    isAnyProjectLinked &&
                    Object.keys(plugin_constants.form.subscriptionModal).map((field) => (
                        <Form
                            key={plugin_constants.form.subscriptionModal[field as SubscriptionModalFields].label}
                            fieldConfig={plugin_constants.form.subscriptionModal[field as SubscriptionModalFields]}
                            value={formFields[field as SubscriptionModalFields]}
                            optionsList={getDropDownOptions(field as SubscriptionModalFields)}
                            onChange={(newValue) => onChangeOfFormField(field as SubscriptionModalFields, newValue)}
                            error={errorState[field as SubscriptionModalFields]}
                            isDisabled={isLoading}
                        />
                    ))
                }
            </>
        </Modal>
    );
};

export default SubscribeModal;
