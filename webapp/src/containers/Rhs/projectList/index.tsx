import React, {useEffect, useState} from 'react';
import {useDispatch} from 'react-redux';

import ProjectCard from 'components/card/project';
import EmptyState from 'components/emptyState';
import LinearLoader from 'components/loader/linear';
import ConfirmationModal from 'components/modal/confirmationModal';

import {setProjectDetails} from 'reducers/projectDetails';
import {showLinkModal} from 'reducers/linkModal';
import usePluginApi from 'hooks/usePluginApi';

import plugin_constants from 'plugin_constants';

const ProjectList = () => {
    const [showConfirmationModal, setShowConfirmationModal] = useState(false);
    const [projectNameToBeUnlinked, setProjectNameToBeUnlinked] = useState<string | null>(null);

    const dispatch = useDispatch();
    const usePlugin = usePluginApi();

    const handleProjectTitleClick = (projectDetails: ProjectDetails) => {
        dispatch(setProjectDetails(projectDetails));
    };

    const handleOpenLinkProjectModal = () => {
        dispatch(showLinkModal([]));
    };

    const handleUnlinkProject = (projectDetails: ProjectDetails) => {
        setProjectNameToBeUnlinked(projectDetails.projectName);
        setShowConfirmationModal(true);
    };

    useEffect(() => {
        usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName);
    }, []);

    return (
        <>
            <p className='rhs-title'>{'Linked Projects'}</p>
            {
                <ConfirmationModal
                    isOpen={showConfirmationModal}
                    onHide={() => setShowConfirmationModal(false)}
                    onConfirm={() => setShowConfirmationModal(false)}
                    confirmBtnText='Unlink'
                    description={`Are you sure you want to unlink ${projectNameToBeUnlinked}?`}
                    title='Confirm Project Unlink'
                />
            }
            {
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isLoading && (
                    <LinearLoader/>
                )
            }
            {
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).isSuccess &&
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).data && (
                    usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).data ?
                        usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName).data?.map((item) => (
                            <ProjectCard
                                onProjectTitleClick={handleProjectTitleClick}
                                projectDetails={item}
                                key={item.projectID}
                                handleUnlinkProject={handleUnlinkProject}
                            />
                        ),
                        ) :
                        <EmptyState
                            title='No Project Linked'
                            subTitle={{text: 'You can link a project by clicking the below button or using the slash command', slashCommand: '/azuredevops link'}}
                            buttonText='Link new project'
                            buttonAction={handleOpenLinkProjectModal}
                        />)
            }
        </>
    );
};

export default ProjectList;
