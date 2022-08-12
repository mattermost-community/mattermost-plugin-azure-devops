import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import ProjectCard from 'components/card/project';
import EmptyState from 'components/emptyState';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import {setProjectDetails} from 'reducers/projectDetails';
import {toggleShowLinkModal, toggleIsLinked} from 'reducers/linkModal';
import {getLinkModalState} from 'selectors';
import usePluginApi from 'hooks/usePluginApi';
import plugin_constants from 'plugin_constants';

const ProjectList = () => {
    // State variables
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [projectToBeUnlinked, setProjectToBeUnlinked] = useState<ProjectDetails>();

    // Hooks
    const dispatch = useDispatch();
    const usePlugin = usePluginApi();

    // Fetch linked projects list
    const fetchLinkedProjectsList = () => usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);

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
    const handleConfirmUnlinkProject = async () => {
        const unlinkProjectStatus = await usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectToBeUnlinked);

        if (unlinkProjectStatus) {
            fetchLinkedProjectsList();
            setShowConfirmationModal(false);
        }
    };

    // Fetch the linked projects list when RHS is opened
    useEffect(() => {
        fetchLinkedProjectsList();
    }, []);

    // Fetch the linked projects list when new project is linked
    useEffect(() => {
        if (getLinkModalState(usePlugin.state).isLinked) {
            dispatch(toggleIsLinked(false));
            fetchLinkedProjectsList();
        }
    }, [getLinkModalState(usePlugin.state).isLinked]);

    const data = usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).data as ProjectDetails[];

    return (
        <>
            <p className='rhs-title'>{'Linked Projects'}</p>
            {
                <ConfirmationModal
                    isOpen={showConfirmationModal}
                    onHide={() => setShowConfirmationModal(false)}
                    onConfirm={handleConfirmUnlinkProject}
                    isLoading={usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.unlinkProject.apiServiceName, projectToBeUnlinked).isLoading}
                    confirmBtnText='Unlink'
                    description={`Are you sure you want to unlink ${projectToBeUnlinked?.projectName}?`}
                    title='Confirm Project Unlink'
                />
            }
            {
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isLoading && (
                    <LinearLoader/>
                )
            }
            {
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isSuccess && (
                    data && data.length > 0 ?
                        <>
                            {
                                data?.map((item) => (
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
                            title='No Project Linked'
                            subTitle={{text: 'Link a project by clicking the button below'}}
                            buttonText='Link new project'
                            buttonAction={handleOpenLinkProjectModal}
                        />)
            }
        </>
    );
};

export default ProjectList;
