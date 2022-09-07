import React, {memo, useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';

import EmptyState from 'components/emptyState';
import SubscriptionCard from 'components/card/subscription';
import BackButton from 'components/buttons/backButton';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';
import ToggleSwitch from 'components/toggleSwitch';
import PrimaryButton from 'components/buttons/primaryButton';

import plugin_constants from 'plugin_constants';

import {resetProjectDetails} from 'reducers/projectDetails';
import {toggleIsSubscribed, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleIsLinkedProjectListChanged} from 'reducers/linkModal';
import {toggleIsSubscriptionDeleted} from 'reducers/websocketEvent';
import {getSubscribeModalState, getWebsocketEventState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {getCurrentChannelSubscriptions} from 'utils/filterData';

const ProjectDetails = memo((projectDetails: ProjectDetails) => {
    // Hooks
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, makeApiRequest, getApiState, state} = usePluginApi();

    // State variables
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();
    const [showAllSubscriptions, setShowAllSubscriptions] = useState(false);
    const [subscriptionList, setSubscriptionList] = useState<SubscriptionDetails[]>();
    const {entities} = useSelector((globalState: GlobalState) => globalState);
    const {currentChannelId} = entities.channels;

    const project: FetchSubscriptionList = {project: projectDetails.projectName};
    const {data, isLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, project);
    const subscriptionListReturnedByApi = data as SubscriptionDetails[];

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

    // Opens subscription modal
    const handleSubscriptionModal = () => {
        dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs: []}));
    };

    // Opens a confirmation modal to confirm unlinking a project
    const handleUnlinkProject = () => {
        setShowProjectConfirmationModal(true);
    };

    // Opens a confirmation modal to confirm deletion of a subscription
    const handleDeleteSubscription = (subscriptionDetails: SubscriptionDetails) => {
        setSubscriptionToBeDeleted({
            organization: subscriptionDetails.organizationName,
            project: subscriptionDetails.projectName,
            eventType: subscriptionDetails.eventType,
            channelID: subscriptionDetails.channelID,
        });
        setShowSubscriptionConfirmationModal(true);
    };

    // Handles unlinking a project and fetching the modified project list
    const handleConfirmUnlinkProject = () => {
        makeApiRequestWithCompletionStatus(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);
    };

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: projectDetails,
        handleSuccess: () => {
            dispatch(toggleIsLinkedProjectListChanged(true));
            handleResetProjectDetails();
            setShowProjectConfirmationModal(false);
        },
    });

    // Fetch subscription list
    const fetchSubscriptionList = () => makeApiRequest(
        plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        project,
    );

    // Handles deletion of a subscription and fetching the modified subscription list
    const handleConfirmDeleteSubscription = () => {
        makeApiRequestWithCompletionStatus(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);
    };

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName,
        payload: subscriptionToBeDeleted,
        handleSuccess: () => {
            fetchSubscriptionList();
            setShowSubscriptionConfirmationModal(false);
        },
    });

    // Handles switch toggle.
    const handleToggle = () => {
        if (showAllSubscriptions) {
            setSubscriptionList(getCurrentChannelSubscriptions(subscriptionListReturnedByApi, currentChannelId));
        } else {
            setSubscriptionList(subscriptionListReturnedByApi);
        }
        setShowAllSubscriptions(!showAllSubscriptions);
    };

    // Reset the state when the component is unmounted
    useEffect(() => {
        fetchSubscriptionList();
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    useEffect(() => {
        // Update subscription list only when it does not match with the current data
        if (subscriptionListReturnedByApi !== subscriptionList) {
            if (showAllSubscriptions) {
                setSubscriptionList(subscriptionListReturnedByApi);
            } else {
                setSubscriptionList(getCurrentChannelSubscriptions(subscriptionListReturnedByApi, currentChannelId));
            }
        }
    }, [subscriptionListReturnedByApi]);

    // Update subscription list on switching channels
    useEffect(() => {
        setShowAllSubscriptions(false);
        setSubscriptionList(getCurrentChannelSubscriptions(subscriptionListReturnedByApi, currentChannelId));
    }, [currentChannelId]);

    // Fetch the subscription list when new subscription is created
    useEffect(() => {
        if (getSubscribeModalState(state).isCreated) {
            dispatch(toggleIsSubscribed(false));
            fetchSubscriptionList();
        }
    }, [getSubscribeModalState(state).isCreated]);

    const {isLoading: isUnlinkProjectLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);
    const {isLoading: isDeleteSubscriptionLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

    // Update the subscription list on RHS when a subscription is deleted using the slash command
    useEffect(() => {
        if (getWebsocketEventState(state).isSubscriptionDeleted) {
            fetchSubscriptionList();
            dispatch(toggleIsSubscriptionDeleted(false));
        }
    }, [getWebsocketEventState(state).isSubscriptionDeleted]);

    return (
        <>
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={isUnlinkProjectLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectDetails.projectName}?`}
                title='Confirm Project Unlink'
            />
            <ConfirmationModal
                isOpen={showSubscriptionConfirmationModal}
                onHide={() => setShowSubscriptionConfirmationModal(false)}
                onConfirm={handleConfirmDeleteSubscription}
                isLoading={isDeleteSubscriptionLoading}
                confirmBtnText='Delete'
                description='Are you sure you want to delete this subscription?'
                title='Confirm Delete Subscription'
            />
            {isLoading && <LinearLoader extraClass='top-0'/>}
            <ToggleSwitch
                active={showAllSubscriptions}
                onChange={handleToggle}
                label={'Show All Subscriptions'}
                labelPositioning='right'
            />
            <div className='d-flex align-item-center margin-bottom-15'>
                <BackButton onClick={handleResetProjectDetails}/>
                <p className='rhs-title'>{projectDetails.projectName}</p>
                <PrimaryButton
                    text='Unlink'
                    iconName='fa fa-chain-broken'
                    extraClass='rhs-project-details-unlink-button'
                    onClick={handleUnlinkProject}
                />
            </div>
            {
                subscriptionList?.length ? (
                    <>
                        {
                            subscriptionList.map((item) => (
                                <SubscriptionCard
                                    subscriptionDetails={item}
                                    key={item.mattermostUserID}
                                    handleDeleteSubscrption={handleDeleteSubscription}
                                />
                            ),
                            )
                        }
                        <div className='rhs-project-list-wrapper'>
                            <button
                                onClick={handleSubscriptionModal}
                                className='plugin-btn no-data__btn btn btn-primary project-list-btn'
                            >
                                {'Add New Subscription'}
                            </button>
                        </div>
                    </>
                ) : (
                    <EmptyState
                        title='No subscriptions yet'
                        subTitle={{text: 'You can add a subscription by clicking the below button.'}}
                        buttonText='Add new subscription'
                        buttonAction={handleSubscriptionModal}
                        icon='subscriptions'
                        wrapperExtraClass='margin-top-80'
                    />
                )
            }
        </>
    );
});

export default ProjectDetails;
