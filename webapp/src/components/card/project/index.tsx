import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';

import {onPressingEnterKey} from 'utils';

type ProjectCardProps = {
    onProjectTitleClick: (projectDetails: ProjectDetails) => void
    handleUnlinkProject: () => void
    projectDetails: ProjectDetails
}

const ProjectCard = ({onProjectTitleClick, projectDetails: {organizationName, projectName}, projectDetails, handleUnlinkProject}: ProjectCardProps) => {
    return (
        <BaseCard>
            <div className='d-flex'>
                <div className='project-details'>
                    <p className='margin-bottom-10'>
                        <span
                            aria-label={projectName}
                            role='button'
                            tabIndex={0}
                            className='font-size-14 font-bold link-title'
                            onKeyDown={() => onPressingEnterKey(event, () => onProjectTitleClick(projectDetails))}
                            onClick={() => onProjectTitleClick(projectDetails)}
                        >
                            {projectName}
                        </span>
                    </p>
                    <p className='font-size-14'>{organizationName}</p>
                </div>
                <div className='button-wrapper'>
                    <IconButton
                        tooltipText='Unlink project'
                        iconClassName='fa fa-chain-broken'
                        extraClass='unlink-button'
                        onClick={handleUnlinkProject}
                    />
                </div>
            </div>
        </BaseCard>
    );
};

export default ProjectCard;

