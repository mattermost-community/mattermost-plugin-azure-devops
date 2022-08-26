import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import SubscriptionCard from 'components/card/subscription';
import IconButton from 'components/buttons/iconButton';
import BackButton from 'components/buttons/backButton';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import usePluginApi from 'hooks/usePluginApi';
import {resetProjectDetails} from 'reducers/projectDetails';
import plugin_constants from 'plugin_constants';
import EmptyState from 'components/emptyState';
import {toggleIsSubscribed, toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {getSubscribeModalState} from 'selectors';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

const ProjectDetails = (projectDetails: ProjectDetails) => {
    // State variables
    const [showProjectConfirmationModal, setShowProjectConfirmationModal] = useState(false);
    const [showSubscriptionConfirmationModal, setShowSubscriptionConfirmationModal] = useState(false);
    const [subscriptionToBeDeleted, setSubscriptionToBeDeleted] = useState<SubscriptionPayload>();

    // Hooks
    const dispatch = useDispatch();
    const usePlugin = usePluginApi();

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
        usePlugin.makeApiRequestWithCompletionStatus(
            plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
            projectDetails,
        );
    };

    const handleAfterUnlinkingProject = () => {
        usePlugin.makeApiRequest(
            plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
        );
        handleResetProjectDetails();
        setShowProjectConfirmationModal(false);
    };

    // Handle sucess/error response of API call made to unlink project
    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: projectDetails,
        handleSuccess: handleAfterUnlinkingProject,
    });

    // Handles deletion of a subscription and fetching the modified subscription list
    const handleConfirmDeleteSubscription = async () => {
        const deleteSubscriptionStatus = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted);

        if (deleteSubscriptionStatus) {
            fetchSubscriptionList();
            setShowSubscriptionConfirmationModal(false);
        }
    };

    const project: FetchSubscriptionList = {project: projectDetails.projectName};
    const {data, isLoading} = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName, project);
    const subscriptionList = data as SubscriptionDetails[];

    // Fetch subscription list
    const fetchSubscriptionList = () => usePlugin.makeApiRequest(
        plugin_constants.pluginApiServiceConfigs.getSubscriptionList.apiServiceName,
        project,
    );

    // Reset the state when the component is unmounted
    useEffect(() => {
        fetchSubscriptionList();
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    // Fetch the subscription list when new subscription is created
    useEffect(() => {
        if (getSubscribeModalState(usePlugin.state).isCreated) {
            dispatch(toggleIsSubscribed(false));
            fetchSubscriptionList();
        }
    }, [getSubscribeModalState(usePlugin.state).isCreated]);

    return (
        <>
            <BackButton onClick={handleResetProjectDetails}/>
            <ConfirmationModal
                isOpen={showProjectConfirmationModal}
                onHide={() => setShowProjectConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails).isLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectDetails.projectName}?`}
                title='Confirm Project Unlink'
            />
            <ConfirmationModal
                isOpen={showSubscriptionConfirmationModal}
                onHide={() => setShowSubscriptionConfirmationModal(false)}
                onConfirm={handleConfirmDeleteSubscription}
                isLoading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.deleteSubscription.apiServiceName, subscriptionToBeDeleted).isLoading}
                confirmBtnText='Delete'
                description='Are you sure you want to delete this subscription ?'
                title='Confirm Delete Subscription'
            />
            {isLoading && <LinearLoader/>}
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
                subscriptionList && subscriptionList.length ?
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
                    </> :
                    <EmptyState
                        title='No subscriptions yet'
                        subTitle={{text: 'You can link a subscription by clicking the below button.'}}
                        buttonText='Add new subscription'
                        buttonAction={handleSubscriptionModal}
                        icon='subscriptions'
                        wrapperExtraClass='margin-top-80'
                    />
            }
        </>
    );
};

export default ProjectDetails;
