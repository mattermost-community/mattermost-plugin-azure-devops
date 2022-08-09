import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';

import {onPressingEnterKey} from 'utils';

import './styles.scss';

type ProjectCardProps = {
    onProjectTitleClick: (projectDetails: ProjectDetails) => void;
    projectDetails: ProjectDetails
}

const ProjectCard = ({onProjectTitleClick, projectDetails: {organization, title}, projectDetails}: ProjectCardProps) => {
    return (
        <BaseCard>
            <div className='d-flex'>
                <div className='project-details'>
                    <p className='margin-bottom-10'>
                        <span
                            aria-label={title}
                            role='button'
                            tabIndex={0}
                            className='font-size-14 font-bold link-title'
                            onKeyDown={() => onPressingEnterKey(event, () => onProjectTitleClick(projectDetails))}
                            onClick={() => onProjectTitleClick(projectDetails)}
                        >
                            {title}
                        </span>
                    </p>
                    <p className='font-size-14'>{organization}</p>
                </div>
                <div className='button-wrapper'>
                    <IconButton
                        tooltipText='Unlink project'
                        iconClassName='fa fa-chain-broken'
                        extraClass='project-details-unlink-button unlink-button'
                    />
                </div>
            </div>
        </BaseCard>
    );
};

export default ProjectCard;

