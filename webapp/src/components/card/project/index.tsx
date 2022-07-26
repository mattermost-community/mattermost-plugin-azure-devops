import React from 'react'

import BaseCard from 'components/card/base'
import IconButton from 'components/buttons/iconButton'

import './styles.scss'

const ProjectCard = ({ title, organization }: Pick<ProjectItem, 'title' | 'organization'>) => {
  return (
    <BaseCard>
      <div className='d-flex'>
        <div className='project-details'>
          <p className='margin-bottom-10 font-size-14 font-bold link-title'>{title}</p>
          <p className='font-size-14'>{organization}</p>
        </div>
        <div className='button-wrapper'>
          <IconButton tooltipText='unlink project' iconClassName='fa fa-trash-o' extraClass='delete-button' iconColor='danger'/>
        </div>
      </div>
    </BaseCard>
  )
}

export default ProjectCard;
