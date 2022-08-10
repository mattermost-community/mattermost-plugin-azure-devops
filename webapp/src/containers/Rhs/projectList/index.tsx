import React from 'react';
import {useDispatch} from 'react-redux';

import ProjectCard from 'components/card/project';

import {setProjectDetails} from 'reducers/projectDetails';

// TODO: dummy data, remove later
const data: ProjectDetails[] = [
    {
        id: 'abc',
        title: 'Project A',
        organization: 'Organization Name',
    },
    {
        id: 'abc1',
        title: 'Project B',
        organization: 'Organization Name',
    },
    {
        id: 'abc2',
        title: 'Project C',
        organization: 'Organization Name',
    },
];

const ProjectList = () => {
    const dispatch = useDispatch();

    const handleProjectTitleClick = (projectDetails: ProjectDetails) => {
        dispatch(setProjectDetails(projectDetails));
    };

    return (
        <>
            <p className='rhs-title'>{'Linked Projects'}</p>
            {
                data.map((item) => (
                    <ProjectCard
                        onProjectTitleClick={handleProjectTitleClick}
                        projectDetails={{id: item.id, title: item.title, organization: item.organization}}
                        key={item.id}
                    />
                ),
                )
            }
        </>
    );
};

export default ProjectList;
