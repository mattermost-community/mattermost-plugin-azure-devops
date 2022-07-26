import React from 'react';

import ProjectCard from 'components/card/project';

// TODO: dummy data, remove later
const data: ProjectItem[] = [
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
    return (
        <>
            <p className='rhs-title'>{'Linked Projects'}</p>
            {
                data.map((item) => (
                    <ProjectCard
                        title={item.title}
                        organization={item.organization}
                        key={item.id}
                    />
                ),
                )
            }
        </>
    );
};

export default ProjectList;
