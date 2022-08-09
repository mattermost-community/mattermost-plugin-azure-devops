import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState} from 'selectors';

import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const usePlugin = usePluginApi();

    return (
        <div className='height-100vh bg-sidebar padding-25'>
            {
                getprojectDetailsState(usePlugin.state).id ?
                    <ProjectDetails title={getprojectDetailsState(usePlugin.state).title}/> :
                    <ProjectList/>
            }
        </div>
    );
};

export default Rhs;
