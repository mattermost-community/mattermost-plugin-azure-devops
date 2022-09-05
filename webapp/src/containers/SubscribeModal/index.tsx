import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';

import Modal from 'components/modal';

import usePluginApi from 'hooks/usePluginApi';
import {getSubscribeModalState} from 'selectors';
import plugin_constants from 'plugin_constants';
import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import Dropdown from 'components/dropdown';
import {getOrganizationList, getProjectList} from 'utils';

import './styles.scss';

const SubscribeModal = () => {
    const eventTypeOptions = [
        {
            value: 'create',
            label: 'Create',
        },
        {
            value: 'update',
            label: 'Update',
        },
        {
            value: 'delete',
            label: 'Delete',
        },
    ];

    const [subscriptionDetails, setSubscriptionDetails] = useState<SubscriptionPayload>({
        organization: '',
        project: '',
        eventType: '',
        channelID: '',
    });
    const [errorState, setErrorState] = useState<SubscriptionPayload>({
        organization: '',
        project: '',
        eventType: '',
        channelID: '',
    });
    const [channelOptions, setChannelOptions] = useState<DropdownOptionType[]>([]);
    const [organizationOptions, setOrganizationOptions] = useState<DropdownOptionType[]>([]);
    const [projectOptions, setProjectOptions] = useState<DropdownOptionType[]>([]);
    const {entities} = useSelector((state: GlobalState) => state);
    const usePlugin = usePluginApi();
    const visibility = getSubscribeModalState(usePlugin.state).visibility;
    const dispatch = useDispatch();

    // Get ProjectList State
    const getProjectState = () => {
        const {isLoading, isSuccess, isError, data} = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        return {isLoading, isSuccess, isError, data: data as ProjectDetails[]};
    };

    // Get ChannelList State
    const getChannelState = () => {
        const {isLoading, isSuccess, isError, data} = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName, {teamId: entities.teams.currentTeamId});
        return {isLoading, isSuccess, isError, data: data as ChannelList[]};
    };

    useEffect(() => {
        if (!getChannelState().data) {
            usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getChannels.apiServiceName, {teamId: entities.teams.currentTeamId});
        }
        if (!getProjectState().data) {
            usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
        }
    }, []);

    useEffect(() => {
        // Pre-select the dropdown value in case of single option.
        if (organizationOptions.length === 1) {
            setSubscriptionDetails((value) => ({...value, organization: organizationOptions[0].value}));
        }
        if (projectOptions.length === 1) {
            setSubscriptionDetails((value) => ({...value, project: projectOptions[0].value}));
        }
        if (channelOptions.length === 1) {
            setSubscriptionDetails((value) => ({...value, channelID: channelOptions[0].value}));
        }
    }, [projectOptions, organizationOptions, channelOptions]);

    // Function to hide the modal and reset all the states.
    const onHide = () => {
        setSubscriptionDetails({
            organization: '',
            project: '',
            eventType: '',
            channelID: '',
        });
        setErrorState({
            organization: '',
            project: '',
            eventType: '',
            channelID: '',
        });
        dispatch(toggleShowSubscribeModal({isVisible: false, commandArgs: []}));
    };

    // Set organization name
    const onOrganizationChange = (value: string) => {
        setErrorState({...errorState, organization: ''});
        setSubscriptionDetails({...subscriptionDetails, organization: value});
    };

    // Set project name
    const onProjectChange = (value: string) => {
        setErrorState({...errorState, project: ''});
        setSubscriptionDetails({...subscriptionDetails, project: value});
    };

    // Set event type
    const onEventTypeChange = (value: string) => {
        setErrorState({...errorState, eventType: ''});
        setSubscriptionDetails({...subscriptionDetails, eventType: value});
    };

    // Set channel name
    const onChannelChange = (value: string) => {
        setErrorState({...errorState, channelID: ''});
        setSubscriptionDetails({...subscriptionDetails, channelID: value});
    };

    // Handles on confirming subscription
    const onConfirm = () => {
        const errorStateChanges: SubscriptionPayload = {
            organization: '',
            project: '',
            channelID: '',
            eventType: '',
        };
        if (subscriptionDetails.organization === '') {
            errorStateChanges.organization = 'Organization is required';
        }
        if (subscriptionDetails.project === '') {
            errorStateChanges.project = 'Project is required';
        }
        if (subscriptionDetails.eventType === '') {
            errorStateChanges.eventType = 'Event type is required';
        }
        if (subscriptionDetails.channelID === '') {
            errorStateChanges.channelID = 'Channel name is required';
        }
        if (errorStateChanges.organization || errorStateChanges.project || errorStateChanges.channelID || errorStateChanges.eventType) {
            setErrorState(errorStateChanges);
            return;
        }

        // Make POST api request
        createSubscription();
    };

    // Make POST api request to create subscription
    const createSubscription = async () => {
        const createSubscriptionResponse = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName, subscriptionDetails);
        if (createSubscriptionResponse) {
            onHide();
        }
    };

    useEffect(() => {
        const channelList = getChannelState().data;
        if (channelList) {
            setChannelOptions(channelList.map((channel) => ({label: <span><i className='fa fa-globe dropdown-option-icon'/>{channel.display_name}</span>, value: channel.id})));
        }

        const projectList = getProjectState().data;
        if (projectList) {
            setProjectOptions(getProjectList(projectList));
            setOrganizationOptions(getOrganizationList(projectList));
        }
    }, [
        getChannelState().isLoading,
        getProjectState().isLoading,
    ]);

    const APIResponse = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.createSubscription.apiServiceName, subscriptionDetails);

    return (
        <Modal
            show={visibility}
            title='Create subscription'
            onHide={onHide}
            onConfirm={onConfirm}
            confirmBtnText='Create subsciption'
            confirmDisabled={APIResponse.isLoading}
            cancelDisabled={APIResponse.isLoading}
            loading={APIResponse.isLoading}
            error={APIResponse.isError ? 'Some error occurred. Please try again later.' : ''}
        >
            <>
                <Dropdown
                    placeholder='Organization name'
                    value={subscriptionDetails.organization}
                    onChange={(newValue) => onOrganizationChange(newValue)}
                    options={organizationOptions}
                    required={true}
                    error={errorState.organization}
                    disabled={APIResponse.isLoading}
                />
                <Dropdown
                    placeholder='Project name'
                    value={subscriptionDetails.project}
                    onChange={(newValue) => onProjectChange(newValue)}
                    options={projectOptions}
                    required={true}
                    error={errorState.project}
                    disabled={APIResponse.isLoading}
                />
                <Dropdown
                    placeholder='Event type'
                    value={subscriptionDetails.eventType}
                    onChange={(newValue) => onEventTypeChange(newValue)}
                    options={eventTypeOptions}
                    required={true}
                    error={errorState.eventType}
                    disabled={APIResponse.isLoading}
                />
                <Dropdown
                    placeholder='Channel name'
                    value={subscriptionDetails.channelID}
                    onChange={(newValue) => onChannelChange(newValue)}
                    options={channelOptions}
                    required={true}
                    error={errorState.channelID}
                    disabled={APIResponse.isLoading}
                />
            </>
        </Modal>
    );
};

export default SubscribeModal;
