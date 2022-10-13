import React, {memo, useEffect, useMemo, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import InfiniteScroll from 'react-infinite-scroll-component';

import {GlobalState} from 'mattermost-redux/types/store';

import EmptyState from 'components/emptyState';
import SubscriptionCard from 'components/card/subscription';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';
import Spinner from 'components/loader/spinner';

import pluginConstants from 'pluginConstants';

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
    const {defaultPage, defaultPerPageLimit, defaultSubscriptionFilters} = pluginConstants.common;

    // State variables
    const [paginationQueryParams, setPaginationQueryParams] = useState<PaginationQueryParams>({
        page: defaultPage,
        per_page: defaultPerPageLimit,
    });
    const [subscriptionList, setSubscriptionList] = useState<SubscriptionDetails[]>([]);
    const [showAllSubscriptions, setShowAllSubscriptions] = useState(false);
    const [filter, setFilter] = useState<SubscriptionFilters>(defaultSubscriptionFilters);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();
    const [deleteConfirmationModalError, setDeleteConfirmationModalError] = useState<ConfirmationModalErrorPanelProps | null>(null);

    // Hooks
    const {currentChannelId} = useSelector((reduxState: GlobalState) => reduxState.entities.channels);
    const {currentTeamId} = useSelector((reduxState: GlobalState) => reduxState.entities.teams);
    const previousState = usePreviousState({currentChannelId});
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState, state} = usePluginApi();

    const subscriptionListApiParams = useMemo<FetchSubscriptionList>(() => ({
        project: projectName,
        channel_id: showAllSubscriptions ? '' : currentChannelId,
        page: paginationQueryParams.page,
        per_page: paginationQueryParams.per_page,
        created_by: filter.createdBy,
        service_type: filter.serviceType,
        event_type: filter.eventType,
        team_id: currentTeamId,
    }), [projectName, currentChannelId, currentTeamId, showAllSubscriptions, paginationQueryParams, filter]);

    const {data, isLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, subscriptionListApiParams);
    const subscriptionListReturnedByApi = data as SubscriptionDetails[] || [];
    const hasMoreSubscriptions = useMemo<boolean>(() => (
        subscriptionListReturnedByApi.length !== 0 && subscriptionListReturnedByApi.length === defaultPerPageLimit
    ), [subscriptionListReturnedByApi]);

    const handlePagination = (reset = false) => {
        if (reset) {
            setSubscriptionList([]);
        }

        setPaginationQueryParams({
            ...paginationQueryParams,
            page: reset ? defaultPage : paginationQueryParams.page + 1,
        });
    };

    const handleSetFilter = (newFilter: SubscriptionFilters) => {
        setFilter(newFilter);
        handlePagination(true);
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
        pluginConstants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        subscriptionListApiParams,
    );

    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        payload: subscriptionListApiParams,
        handleSuccess: () => {
            setSubscriptionList([...subscriptionList, ...subscriptionListReturnedByApi]);
        },
    });

    // Handles deletion of a subscription and fetching the modified subscription list
    const handleConfirmDeleteSubscription = () => makeApiRequestWithCompletionStatus(pluginConstants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.deleteSubscription.apiServiceName,
        payload: subscriptionToBeDeleted,
        handleSuccess: () => {
            handlePagination(true);
            setShowSubscriptionConfirmationModal(false);
        },
        handleError: (error) => {
            const errorMessage = utils.getErrorMessage(true, 'ConfirmationModal', error);
            if (errorMessage === pluginConstants.messages.error.subscriptionNotFound) {
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
        if (previousState?.currentChannelId !== currentChannelId) {
            if (showAllSubscriptions) {
                return;
            }
            setSubscriptionList([]);
        }

        fetchSubscriptionList();
    }, [
        currentChannelId,
        currentTeamId,
        paginationQueryParams,
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

    const {isLoading: isDeleteSubscriptionLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

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
                setFilter={handleSetFilter}
                setSubscriptionList={setSubscriptionList}
            />
            {
                subscriptionList.length ? (
                    <InfiniteScroll
                        dataLength={defaultPerPageLimit}
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
