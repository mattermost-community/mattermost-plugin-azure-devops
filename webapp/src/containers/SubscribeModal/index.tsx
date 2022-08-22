import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';

import Modal from 'components/modal';
import Form from 'components/form';

import plugin_constants from 'plugin_constants';

import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import usePluginApi from 'hooks/usePluginApi';
import useForm from 'hooks/useForm';

import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {getSubscribeModalState} from 'selectors';

import Utils, {getOrganizationList, getProjectList} from 'utils';

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
    const {getApiState, makeApiRequest, makeApiRequestWithCompletionStatus, state} = usePluginApi();
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

    // Get ProjectList State
    const getProjectState = () => {
        const {isLoading, isSuccess, isError, data} = getApiState(
            plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
        );
        return {isLoading, isSuccess, isError, data: data as ProjectDetails[]};
    };

    // Get ChannelList State
    const getChannelState = () => {
        const {isLoading, isSuccess, isError, data} = getApiState(
            plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName,
            {teamId: entities.teams.currentTeamId},
        );
        return {isLoading, isSuccess, isError, data: data as ChannelList[]};
    };

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

    // Handles on confirming create subscription
    const onConfirm = () => {
        if (!isErrorInFormValidation()) {
            // Make POST api request to create subscription
            makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName,
                formFields,
            );
        }
    };

    // Observe for the change in redux state after API call and do the required actions
    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName,
        handleSuccess: () => resetModalState(true),
        payload: formFields,
    });

    // Make API request to fetch channel and project list
    useEffect(() => {
        if (!getChannelState().data) {
            makeApiRequest(
                plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName,
                {teamId: entities.teams.currentTeamId},
            );
        }
        if (!getProjectState().data) {
            makeApiRequest(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        }
    }, [visibility]);

    // Pre-select the dropdown value in case of single option.
    useEffect(() => {
        if (organizationOptions.length === 1) {
            setSpecificFieldValue('organization', organizationOptions[0].value);
        }
        if (projectOptions.length === 1) {
            setSpecificFieldValue('project', projectOptions[0].value);
        }
        if (channelOptions.length === 1) {
            setSpecificFieldValue('channelID', channelOptions[0].value);
        }
    }, [projectOptions, organizationOptions, channelOptions]);

    // Set channel and project list values
    useEffect(() => {
        const channelList = getChannelState().data;
        if (channelList) {
            setChannelOptions(channelList.map((channel) => ({
                label: <span><i className='fa fa-globe dropdown-option-icon'/>{channel.display_name}</span>,
                value: channel.id,
            })));
        }

        const projectList = getProjectState().data;
        if (projectList) {
            setProjectOptions(getProjectList(projectList));
            setOrganizationOptions(getOrganizationList(projectList));
        }
    }, [state]);

    const {isLoading, isError, error} = getApiState(plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName, formFields);

    return (
        <Modal
            show={visibility}
            title='Create subscription'
            onHide={resetModalState}
            onConfirm={onConfirm}
            confirmBtnText='Create subscription'
            confirmDisabled={isLoading}
            cancelDisabled={isLoading}
            loading={isLoading}
            error={Utils.getErrorMessage(isError, 'SubscribeModal', error as unknown as ApiErrorResponse)}
        >
            <>
                {
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
