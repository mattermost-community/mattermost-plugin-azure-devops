import React, {useEffect, useMemo, useState} from 'react';
import {useDispatch} from 'react-redux';

import ProjectCard from 'components/card/project';
import EmptyState from 'components/emptyState';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import pluginConstants from 'pluginConstants';

import {setProjectDetails} from 'reducers/projectDetails';
import {toggleIsLinkedProjectListChanged, toggleShowLinkModal} from 'reducers/linkModal';
import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import utils, {sortProjectList} from 'utils';

const ProjectList = () => {
    // State variables
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [projectToBeUnlinked, setProjectToBeUnlinked] = useState<ProjectDetails>();
    const [deleteSubscriptions, setDeleteSubscriptions] = useState(false);
    const [showDeleteSubscriptionsCheckbox, setShowDeleteSubscriptionsCheckbox] = useState(true);
    const [confirmationModalDescription, setConfirmationModalDescription] = useState(`Are you sure you want to unlink ${projectToBeUnlinked?.projectName}?`);
    const [unlinkConfirmationModalError, setUnlinkConfirmationModalError] = useState<ConfirmationModalErrorPanelProps | null>(null);

    // Hooks
    const dispatch = useDispatch();
    const {getApiState, makeApiRequestWithCompletionStatus} = usePluginApi();

    // Navigates to project details view
    const handleProjectTitleClick = (projectDetails: ProjectDetails) => {
        dispatch(setProjectDetails(projectDetails));
    };

    // Opens link project modal
    const handleOpenLinkProjectModal = () => {
        dispatch(toggleShowLinkModal({isVisible: true, commandArgs: []}));
    };

    /**
     * Opens a confirmation modal to confirm unlinking a project
     * @param projectDetails
     */
    const handleUnlinkProject = (projectDetails: ProjectDetails) => {
        setProjectToBeUnlinked(projectDetails);
        setShowConfirmationModal(true);
    };

    // Handles unlinking a project and fetching the modified project list
    const handleConfirmUnlinkProject = () => {
        makeApiRequestWithCompletionStatus(
            pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
            {
                ...projectToBeUnlinked,
                deleteSubscriptions,
            } as APIRequestPayload,
        );
    };

    // Fetch updated project list and close the unlink confirmation modal
    const handleActionsAfterUnlinkingProject = () => {
        dispatch(toggleIsLinkedProjectListChanged(true));
        setShowConfirmationModal(false);
    };

    // Update the modal when project unlinking fails
    const handleActionsAfterUnlinkingProjectFailed = (err: ApiErrorResponse) => {
        const errorMessage = utils.getErrorMessage(true, 'ConfirmationModal', err);
        if (errorMessage === pluginConstants.messages.error.adminAccessError) {
            setConfirmationModalDescription(pluginConstants.messages.error.adminAccessErrorForUnlinking);
        }

        setUnlinkConfirmationModalError({
            title: errorMessage,
            onSecondaryBtnClick: () => {
                setShowConfirmationModal(false);
                setUnlinkConfirmationModalError(null);
                setShowDeleteSubscriptionsCheckbox(true);
            },
        });

        setDeleteSubscriptions(false);
        setShowDeleteSubscriptionsCheckbox(false);
    };

    // Handle sucess/error response of API call made to unlink project
    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: {...projectToBeUnlinked, deleteSubscriptions} as APIRequestPayload,
        handleSuccess: handleActionsAfterUnlinkingProject,
        handleError: handleActionsAfterUnlinkingProjectFailed,
    });

    const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setDeleteSubscriptions(e.target.checked);
    };

    const deleteSubscriptionsCheckbox = (
        <div>
            <input
                type='checkbox'
                id='deleteSubscriptions'
                onChange={handleCheckboxChange}
            />
            <label>{'Delete all your subscriptions associated with this project'}</label>
        </div>
    );

    const {data, isSuccess, isLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
    const projectsList = data as ProjectDetails[] ?? [];
    const sortedProjectList = useMemo(() => [...projectsList].sort(sortProjectList), [projectsList]); // TODO: Look for best optimisation method here

    return (
        <>
            <p className='rhs-title margin-bottom-15'>{'Linked Projects'}</p>
            {
                <ConfirmationModal
                    isOpen={showConfirmationModal}
                    onHide={() => setShowConfirmationModal(false)}
                    onConfirm={handleConfirmUnlinkProject}
                    isLoading={getApiState(pluginConstants.pluginApiServiceConfigs.unlinkProject.apiServiceName, {...projectToBeUnlinked, deleteSubscriptions} as APIRequestPayload).isLoading}
                    confirmBtnText='Unlink'
                    description={confirmationModalDescription}
                    title='Confirm Project Unlink'
                    showErrorPanel={unlinkConfirmationModalError}
                >
                    {showDeleteSubscriptionsCheckbox ? deleteSubscriptionsCheckbox : <></>}
                </ConfirmationModal>
            }
            {isLoading && <LinearLoader/>}
            {
                isSuccess && (
                    sortedProjectList.length > 0 ?
                        <>
                            {
                                sortedProjectList.map((item: ProjectDetails) => (
                                    <ProjectCard
                                        onProjectTitleClick={handleProjectTitleClick}
                                        projectDetails={item}
                                        key={item.projectID}
                                        handleUnlinkProject={handleUnlinkProject}
                                    />
                                ),
                                )
                            }
                            <div className='rhs-project-list-wrapper'>
                                <button
                                    onClick={handleOpenLinkProjectModal}
                                    className='plugin-btn no-data__btn btn btn-primary project-list-btn'
                                >
                                    {'Link New Project'}
                                </button>
                            </div>
                        </> :
                        <EmptyState
                            title='No project linked' // TODO: create constants for these texts labels/messages
                            subTitle={{text: 'You can link a project by clicking the below button.'}}
                            buttonText='Link New Project'
                            buttonAction={handleOpenLinkProjectModal}
                            wrapperExtraClass='margin-top-80'
                            isLoading={isLoading}
                        />)
            }
        </>
    );
};

export default ProjectList;
