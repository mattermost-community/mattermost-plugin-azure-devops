import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState, getRhsState} from 'selectors';

import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const usePlugin = usePluginApi();

    if (!getRhsState(usePlugin.state).isSidebarOpen) {
        return <></>;
    }

    return (
        <div className='height-rhs bg-sidebar padding-25'>
            {
                getprojectDetailsState(usePlugin.state).projectID ?
                    <ProjectDetails title={getprojectDetailsState(usePlugin.state).projectName}/> :
                    <ProjectList/>
            }
        </div>
    );
};

export default Rhs;
