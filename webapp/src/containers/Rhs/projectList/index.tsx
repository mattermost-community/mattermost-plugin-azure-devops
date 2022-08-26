import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import ProjectCard from 'components/card/project';
import EmptyState from 'components/emptyState';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import {setProjectDetails} from 'reducers/projectDetails';
import {toggleShowLinkModal} from 'reducers/linkModal';
import usePluginApi from 'hooks/usePluginApi';
import plugin_constants from 'plugin_constants';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

const ProjectList = () => {
    // State variables
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [projectToBeUnlinked, setProjectToBeUnlinked] = useState<ProjectDetails>();

    // Hooks
    const dispatch = useDispatch();
    const {getApiState, makeApiRequest, makeApiRequestWithCompletionStatus} = usePluginApi();

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
            plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
            projectToBeUnlinked,
        );
    };

    // Fetch updated project list and close the unlink confirmation modal
    const handleActionsAfterUnlinkingProject = () => {
        makeApiRequest(
            plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
        );
        setShowConfirmationModal(false);
    };

    // Handle sucess/error response of API call made to unlink project
    useApiRequestCompletionState({
        serviceName: plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName,
        payload: projectToBeUnlinked,
        handleSuccess: handleActionsAfterUnlinkingProject,
    });

    const {data, isSuccess, isLoading} = getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
    const projectsList = data as ProjectDetails[];

    return (
        <>
            <p className='rhs-title'>{'Linked Projects'}</p>
            {
                <ConfirmationModal
                    isOpen={showConfirmationModal}
                    onHide={() => setShowConfirmationModal(false)}
                    onConfirm={handleConfirmUnlinkProject}
                    isLoading={getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectToBeUnlinked).isLoading}
                    confirmBtnText='Unlink'
                    description={`Are you sure you want to unlink ${projectToBeUnlinked?.projectName}?`}
                    title='Confirm Project Unlink'
                />
            }
            {isLoading && <LinearLoader/> }
            {
                isSuccess && (
                    projectsList && projectsList.length > 0 ?
                        <>
                            {
                                projectsList?.map((item: ProjectDetails) => (
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
                                    {'Link new project'}
                                </button>
                            </div>
                        </> :
                        <EmptyState
                            title='No project linked'
                            subTitle={{text: 'You can link a project by clicking the below button.'}}
                            buttonText='Link new project'
                            buttonAction={handleOpenLinkProjectModal}
                            wrapperExtraClass='margin-top-80'
                        />)
            }
        </>
    );
};

export default ProjectList;
