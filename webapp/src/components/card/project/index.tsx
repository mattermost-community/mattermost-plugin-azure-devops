import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';

import {onPressingEnterKey} from 'utils';

type ProjectCardProps = {
    onProjectTitleClick: (projectDetails: ProjectDetails) => void
    handleUnlinkProject: (projectDetails: ProjectDetails) => void
    projectDetails: ProjectDetails
}

const ProjectCard = ({onProjectTitleClick, projectDetails: {organizationName, projectName}, projectDetails, handleUnlinkProject}: ProjectCardProps) => {
    return (
        <BaseCard>
            <div className='d-flex'>
                <div className='project-details'>
                    <p
                        aria-label={projectName}
                        role='button'
                        tabIndex={0}
                        className='font-size-14 font-bold link-title text-truncate margin-bottom-10'
                        onKeyDown={() => onPressingEnterKey(event, () => onProjectTitleClick(projectDetails))}
                        onClick={() => onProjectTitleClick(projectDetails)}
                    >
                        {projectName}
                    </p>
                    <p className='font-size-14 text-truncate'>{organizationName}</p>
                </div>
                <div className='button-wrapper'>
                    <IconButton
                        tooltipText='Unlink project'
                        iconClassName='fa fa-chain-broken'
                        extraClass='unlink-button'
                        onClick={() => handleUnlinkProject(projectDetails)}
                    />
                </div>
            </div>
        </BaseCard>
    );
};

export default ProjectCard;

