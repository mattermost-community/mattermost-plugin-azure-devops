import React, {memo, useEffect, useMemo, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import InfiniteScroll from 'react-infinite-scroll-component';

import {GlobalState} from 'mattermost-redux/types/store';

import EmptyState from 'components/emptyState';
import SubscriptionCard from 'components/card/subscription';
import BackButton from 'components/buttons/backButton';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';
import ToggleSwitch from 'components/toggleSwitch';
import PrimaryButton from 'components/buttons/primaryButton';
import Spinner from 'components/loader/spinner';

import plugin_constants from 'plugin_constants';

import {resetProjectDetails} from 'reducers/projectDetails';
import {toggleIsSubscribed, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleIsLinkedProjectListChanged} from 'reducers/linkModal';
import {toggleIsSubscriptionDeleted} from 'reducers/websocketEvent';
import {getSubscribeModalState, getWebsocketEventState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {getIncrementedPaginationQueryParamOffset} from 'utils';
import usePreviousState from 'hooks/usePreviousState';

const ProjectDetails = memo((projectDetails: ProjectDetails) => {
    const {projectName, organizationName} = projectDetails;

    // State variables
    const [paginationQueryParams, setPaginationQueryParams] = useState<PaginationQueryParams>({
        offset: plugin_constants.common.defaultPageOffset,
        limit: plugin_constants.common.defaultPageLimit,
    });
    const [subscriptionList, setSubscriptionList] = useState<SubscriptionDetails[]>([]);
    const [showAllSubscriptions, setShowAllSubscriptions] = useState(false);
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();

    // Hooks
    const {currentChannelId} = useSelector((reduxState: GlobalState) => reduxState.entities.channels);
    const previousState = usePreviousState({currentChannelId});
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState, state} = usePluginApi();

    const subscriptionListApiParams = useMemo<FetchSubscriptionList>(() => ({
        project: projectName,
        channel_id: showAllSubscriptions ? '' : currentChannelId,
        offset: paginationQueryParams.offset,
        limit: paginationQueryParams.limit,
    }), [projectName, currentChannelId, showAllSubscriptions, paginationQueryParams]);

    const {data, isLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, subscriptionListApiParams);
    const subscriptionListReturnedByApi = data as SubscriptionDetails[] || [];
    const hasMoreSubscriptions = useMemo<boolean>(() => (
        subscriptionListReturnedByApi?.length !== 0 && subscriptionListReturnedByApi?.length === plugin_constants.common.defaultPageLimit
    ), [subscriptionListReturnedByApi]);

    const handlePagination = (reset = false) => {
        if (reset) {
            setSubscriptionList([]);
        }

        const {offset} = getIncrementedPaginationQueryParamOffset(paginationQueryParams.offset);
        setPaginationQueryParams({
            ...paginationQueryParams,
            offset: reset ? plugin_constants.common.defaultPageOffset : offset,
        });
    };

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

    // Opens subscription modal
    const handleSubscriptionModal = () => {
        dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs: [], args: [organizationName, projectName]}));
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
    const fetchSubscriptionList = () => makeApiRequestWithCompletionStatus(
        plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        subscriptionListApiParams,
    );

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        payload: subscriptionListApiParams,
        handleSuccess: () => {
            setSubscriptionList([...subscriptionList, ...subscriptionListReturnedByApi]);
        },
    });

    // Handles deletion of a subscription and fetching the modified subscription list
    const handleConfirmDeleteSubscription = () => makeApiRequestWithCompletionStatus(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName,
        payload: subscriptionToBeDeleted,
        handleSuccess: () => {
            handlePagination(true);
            setShowSubscriptionConfirmationModal(false);
        },
    });

    // Reset the state when the component is unmounted
    useEffect(() => {
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    useEffect(() => {
        /**
         * If all subscriptions for a project are already loaded then do not make API calls on switching channel
         */
        if (previousState?.currentChannelId !== currentChannelId && showAllSubscriptions) {
            return;
        }
        fetchSubscriptionList();
    }, [subscriptionListApiParams.channel_id, subscriptionListApiParams.project, subscriptionListApiParams.offset]);

    // Fetch the subscription list when new subscription is created
    useEffect(() => {
        if (getSubscribeModalState(state).isCreated) {
            dispatch(toggleIsSubscribed(false));
            handlePagination(true);
        }
    }, [getSubscribeModalState(state).isCreated]);

    // Update the subscription list on RHS when a subscription is deleted using the slash command
    useEffect(() => {
        if (getWebsocketEventState(state).isSubscriptionDeleted) {
            handlePagination(true);
            dispatch(toggleIsSubscriptionDeleted(false));
        }
    }, [getWebsocketEventState(state).isSubscriptionDeleted]);

    useEffect(() => {
        if (!showAllSubscriptions && subscriptionList?.length) {
            handlePagination(true);
        }
    }, [currentChannelId]);

    const {isLoading: isUnlinkProjectLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);
    const {isLoading: isDeleteSubscriptionLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

    return (
        <>
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={isUnlinkProjectLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectName}?`}
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
                onChange={(active) => {
                    handlePagination(true);
                    setShowAllSubscriptions(active);
                }}
                label={'Show All Subscriptions'}
                labelPositioning='right'
            />
            <div className='d-flex align-item-center margin-bottom-15'>
                <BackButton onClick={handleResetProjectDetails}/>
                <p className='rhs-title'>{projectName}</p>
                <PrimaryButton
                    text='Unlink'
                    iconName='fa fa-chain-broken'
                    extraClass='rhs-project-details-unlink-button'
                    onClick={handleUnlinkProject}
                />
            </div>
            {
                subscriptionList?.length ? (
                    <InfiniteScroll
                        dataLength={plugin_constants.common.defaultPageLimit}
                        next={handlePagination}
                        hasMore={hasMoreSubscriptions}
                        loader={<Spinner/>}
                        endMessage={
                            <p style={{textAlign: 'center'}}>
                                <b>{'You have seen it all'}</b>
                            </p>
                        }
                        scrollableTarget='scrollableArea'
                    >
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
                    </InfiniteScroll>
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
