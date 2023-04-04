import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';

import SVGWrapper from 'components/svgWrapper';
import pluginConstants from 'pluginConstants';
import LabelValuePair from 'components/labelValuePair';
import Tooltip from 'components/tooltip';

type ProjectCardProps = {
    onProjectTitleClick: (projectDetails: ProjectDetails) => void
    handleUnlinkProject: (projectDetails: ProjectDetails) => void
    projectDetails: ProjectDetails
}

const ProjectCard = ({onProjectTitleClick, projectDetails: {organizationName, projectName}, projectDetails, handleUnlinkProject}: ProjectCardProps) => (
    <BaseCard>
        <div className='d-flex'>
            <div className='project-details'>
                <LabelValuePair
                    label={
                        <Tooltip tooltipContent='Project'>
                            <i aria-hidden='true'>
                                <SVGWrapper
                                    width={16}
                                    height={16}
                                    viewBox='0 0 14 12'
                                    className='svg-icon-stroke-color'
                                >
                                    {pluginConstants.SVGIcons.project}
                                </SVGWrapper>
                            </i>
                        </Tooltip>
                    }
                    onClickValue={() => onProjectTitleClick(projectDetails)}
                    value={projectName}
                    labelExtraClassName='margin-top-1'
                />
                <LabelValuePair
                    label={
                        <Tooltip tooltipContent='Organization'>
                            <i aria-hidden='true'>
                                <SVGWrapper
                                    width={13}
                                    height={13}
                                    viewBox='0 0 10 10'
                                >
                                    {pluginConstants.SVGIcons.organization}
                                </SVGWrapper>
                            </i>
                        </Tooltip>
                    }
                    value={organizationName}
                />
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

export default ProjectCard;

