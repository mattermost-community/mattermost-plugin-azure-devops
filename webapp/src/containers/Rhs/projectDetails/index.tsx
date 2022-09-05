import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';

import SubscriptionCard from 'components/card/subscription';
import IconButton from 'components/buttons/iconButton';
import BackButton from 'components/buttons/backButton';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import usePluginApi from 'hooks/usePluginApi';
import { resetProjectDetails } from 'reducers/projectDetails';
import plugin_constants from 'plugin_constants';
import EmptyState from 'components/emptyState';
import { toggleIsSubscribed, toggleShowSubscribeModal } from 'reducers/subscribeModal';
import { getSubscribeModalState } from 'selectors';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import { toggleIsLinkedProjectListChanged } from 'reducers/linkModal';

const ProjectDetails = (projectDetails: ProjectDetails) => {
    // State variables
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();

    // Hooks
    const dispatch = useDispatch();
    const { makeApiRequestWithCompletionStatus, makeApiRequest, getApiState, state } = usePluginApi();

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

    // Opens subscription modal
    const handleSubscriptionModal = () => {
        dispatch(toggleShowSubscribeModal({ isVisible: true, commandArgs: [] }));
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

    // Reset the state when the component is unmounted
    useEffect(() => {
        fetchSubscriptionList();
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    // Fetch the subscription list when new subscription is created
    useEffect(() => {
        if (getSubscribeModalState(state).isCreated) {
            dispatch(toggleIsSubscribed(false));
            fetchSubscriptionList();
        }
    }, [getSubscribeModalState(state).isCreated]);

    const project: FetchSubscriptionList = { project: projectDetails.projectName };
    const { data, isLoading } = getApiState(plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, project);
    const subscriptionList = data as SubscriptionDetails[];

    return (
        <>
            <BackButton onClick={handleResetProjectDetails} />
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails).isLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectDetails.projectName}?`}
                title='Confirm Project Unlink'
            />
            <ConfirmationModal
                isOpen={showSubscriptionConfirmationModal}
                onHide={() => setShowSubscriptionConfirmationModal(false)}
                onConfirm={handleConfirmDeleteSubscription}
                isLoading={getApiState(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted).isLoading}
                confirmBtnText='Delete'
                description='Are you sure you want to delete this subscription ?'
                title='Confirm Delete Subscription'
            />
            {isLoading && <LinearLoader />}
            <div className='d-flex'>
                <p className='rhs-title'>{projectDetails.projectName}</p>
                <IconButton
                    tooltipText='Unlink project'
                    iconClassName='fa fa-chain-broken'
                    extraClass='project-details-unlink-button unlink-button'
                    onClick={() => handleUnlinkProject()}
                />
            </div>
            {
                subscriptionList?.length ? (
                    <>
                        <div className='bottom-divider'>
                            <p className='font-size-14 font-bold margin-0 show-selected'>{'Subscriptions'}</p>
                        </div>
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
                                {'Add new subscription'}
                            </button>
                        </div>
                    </>
                ) : (
                    <EmptyState
                        title='No subscriptions yet'
                        subTitle={{ text: 'You can link a subscription by clicking the below button.' }}
                        buttonText='Add new subscription'
                        buttonAction={handleSubscriptionModal}
                        icon='subscriptions'
                        wrapperExtraClass='margin-top-80'
                    />
                )
            }
        </>
    );
};

export default ProjectDetails;
