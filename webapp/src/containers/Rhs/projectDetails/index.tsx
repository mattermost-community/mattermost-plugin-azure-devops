import React, {memo, useEffect, useMemo, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import InfiniteScroll from 'react-infinite-scroll-component';

import {GlobalState} from 'mattermost-redux/types/store';

import EmptyState from 'components/emptyState';
import SubscriptionCard from 'components/card/subscription';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';
import Spinner from 'components/loader/spinner';

import plugin_constants from 'plugin_constants';

import {toggleIsSubscribed, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleIsSubscriptionDeleted} from 'reducers/websocketEvent';
import {resetProjectDetails} from 'reducers/projectDetails';
import {getSubscribeModalState, getWebsocketEventState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import usePreviousState from 'hooks/usePreviousState';

import utils from 'utils';

import Header from './header';

const ProjectDetails = memo((projectDetails: ProjectDetails) => {
    const {projectName, organizationName} = projectDetails;

    // State variables
    const [paginationQueryParams, setPaginationQueryParams] = useState<PaginationQueryParams>({
        page: plugin_constants.common.defaultPage,
        per_page: plugin_constants.common.defaultPerPageLimit,
    });
    const [subscriptionList, setSubscriptionList] = useState<SubscriptionDetails[]>([]);
    const [showAllSubscriptions, setShowAllSubscriptions] = useState(false);
    const [filter, setFilter] = useState(plugin_constants.common.SubscriptionFilterCreatedBy.me);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();
    const [deleteConfirmationModalError, setDeleteConfirmationModalError] = useState<ConfirmationModalErrorPanel | null>(null);

    // Hooks
    const {currentChannelId} = useSelector((reduxState: GlobalState) => reduxState.entities.channels);
    const previousState = usePreviousState({currentChannelId});
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState, state} = usePluginApi();

    const subscriptionListApiParams = useMemo<FetchSubscriptionList>(() => ({
        project: projectName,
        channel_id: showAllSubscriptions ? '' : currentChannelId,
        page: paginationQueryParams.page,
        per_page: paginationQueryParams.per_page,
        created_by: filter,
    }), [projectName, currentChannelId, showAllSubscriptions, paginationQueryParams, filter]);

    const {data, isLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, subscriptionListApiParams);
    const subscriptionListReturnedByApi = data as SubscriptionDetails[] || [];
    const hasMoreSubscriptions = useMemo<boolean>(() => (
        subscriptionListReturnedByApi.length !== 0 && subscriptionListReturnedByApi.length === plugin_constants.common.defaultPerPageLimit
    ), [subscriptionListReturnedByApi]);

    const handlePagination = (reset = false, fetchList = true) => {
        if (reset) {
            setSubscriptionList([]);
        }
        if (reset && fetchList && paginationQueryParams.page === plugin_constants.common.defaultPage) {
            fetchSubscriptionList();
            return;
        }

        setPaginationQueryParams({
            ...paginationQueryParams,
            page: reset ? plugin_constants.common.defaultPage : paginationQueryParams.page + 1,
        });
    };

    // Opens subscription modal
    const handleSubscriptionModal = () => {
        dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs: [], args: [organizationName, projectName]}));
    };

    // Opens a confirmation modal to confirm deletion of a subscription
    const handleDeleteSubscription = (subscriptionDetails: SubscriptionDetails) => {
        setSubscriptionToBeDeleted({
            organization: subscriptionDetails.organizationName,
            project: subscriptionDetails.projectName,
            eventType: subscriptionDetails.eventType,
            serviceType: subscriptionDetails.serviceType,
            channelID: subscriptionDetails.channelID,
            mmUserID: subscriptionDetails.mattermostUserID,
        });
        setDeleteConfirmationModalError(null);
        setShowSubscriptionConfirmationModal(true);
    };

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
        handleError: (error) => {
            const errorMessage = utils.getErrorMessage(true, 'ConfirmationModal', error);
            if (errorMessage === plugin_constants.messages.error.subscriptionNotFound) {
                handlePagination(true);
                setShowSubscriptionConfirmationModal(false);
                return;
            }
            setDeleteConfirmationModalError({
                title: errorMessage,
                onSecondaryBtnClick: () => setShowSubscriptionConfirmationModal(false),
            });
        },
    });

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

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
    }, [
        subscriptionListApiParams.channel_id,
        subscriptionListApiParams.project,
        subscriptionListApiParams.page,
        subscriptionListApiParams.created_by,
    ]);

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
        if (!showAllSubscriptions && subscriptionList.length) {
            handlePagination(true, false);
        }
    }, [currentChannelId]);

    const {isLoading: isDeleteSubscriptionLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

    return (
        <>
            <ConfirmationModal
                isOpen={showSubscriptionConfirmationModal}
                onHide={() => setShowSubscriptionConfirmationModal(false)}
                onConfirm={handleConfirmDeleteSubscription}
                isLoading={isDeleteSubscriptionLoading}
                confirmBtnText='Delete'
                description='Are you sure you want to delete this subscription?'
                title='Confirm Delete Subscription'
                showErrorPanel={deleteConfirmationModalError}
            />
            {isLoading && <LinearLoader extraClass='top-0'/>}
            <Header
                projectDetails={projectDetails}
                handleResetProjectDetails={handleResetProjectDetails}
                showAllSubscriptions={showAllSubscriptions}
                setShowAllSubscriptions={setShowAllSubscriptions}
                handlePagination={handlePagination}
                filter={filter}
                setFilter={setFilter}
                setSubscriptionList={setSubscriptionList}
            />
            {
                subscriptionList.length ? (
                    <InfiniteScroll
                        dataLength={plugin_constants.common.defaultPerPageLimit}
                        next={handlePagination}
                        hasMore={hasMoreSubscriptions}
                        loader={<Spinner/>}
                        endMessage={
                            <p style={{textAlign: 'center'}}>
                                <b>{'No more subscriptions present.'}</b>
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
                        isLoading={isLoading}
                    />
                )
            }
        </>
    );
});

export default ProjectDetails;
