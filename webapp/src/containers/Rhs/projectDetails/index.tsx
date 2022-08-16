import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import SubscriptionCard from 'components/card/subscription';
import IconButton from 'components/buttons/iconButton';
import BackButton from 'components/buttons/backButton';
import ConfirmationModal from 'components/modal/confirmationModal';

import usePluginApi from 'hooks/usePluginApi';
import {resetProjectDetails} from 'reducers/projectDetails';

import plugin_constants from 'plugin_constants';

// TODO: dummy data, remove later
const data: SubscriptionDetails[] = [
    {
        id: 'abc',
        name: 'Listen for all new tasks created',
        eventType: 'create',
    },
    {
        id: 'abc1',
        name: 'Listen for any task updated',
        eventType: 'update',
    },
    {
        id: 'abc2',
        name: 'Listen for all any task deleted',
        eventType: 'delete',
    },
];

const ProjectDetails = (projectDetails: ProjectDetails) => {
    // State variables
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);

    // Hooks
    const dispatch = useDispatch();
    const usePlugin = usePluginApi();

    const handleResetProjectDetails = () => {
        dispatch(resetProjectDetails());
    };

    /**
     * Opens a confirmation modal to confirm unlinking a project
     */
    const handleUnlinkProject = () => {
        setShowConfirmationModal(true);
    };

    // Handles unlinking a project
    const handleConfirmUnlinkProject = async () => {
        const unlinkProjectStatus = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails);

        if (unlinkProjectStatus) {
            handleResetProjectDetails();
            setShowConfirmationModal(false);
        }
    };

    // Reset the state when the component is unmounted
    useEffect(() => {
        return () => {
            handleResetProjectDetails();
        };
    }, []);

    return (
        <>
            <ConfirmationModal
                isOpen={showConfirmationModal}
                onHide={() => setShowConfirmationModal(false)}
                onConfirm={handleConfirmUnlinkProject}
                isLoading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectDetails).isLoading}
                confirmBtnText='Unlink'
                description={`Are you sure you want to unlink ${projectDetails?.projectName}?`}
                title='Confirm Project Unlink'
            />
            <BackButton onClick={handleResetProjectDetails}/>
            <div className='d-flex'>
                <p className='rhs-title'>{projectDetails.projectName}</p>
                <IconButton
                    tooltipText='Unlink project'
                    iconClassName='fa fa-chain-broken'
                    extraClass='project-details-unlink-button unlink-button'
                    onClick={() => handleUnlinkProject()}
                />
            </div>
            <div className='bottom-divider'>
                <p className='font-size-14 font-bold margin-0 show-selected'>{'Subscriptions'}</p>
            </div>
            {
                data.map((item) => (
                    <SubscriptionCard
                        subscriptionDetails={item}
                        key={item.id}
                    />
                ),
                )
            }
        </>
    );
};

export default ProjectDetails;
