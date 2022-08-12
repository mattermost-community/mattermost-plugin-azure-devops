import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState, getRhsState} from 'selectors';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const usePlugin = usePluginApi();

    if (!getRhsState(usePlugin.state).isSidebarOpen) {
        return <></>;
    }

    return (
        <div className='overflow-auto height-rhs bg-sidebar padding-25'>
            {
                !usePlugin.isUserAccountConnected() && <AccountNotLinked/>
            }
            {
                usePlugin.isUserAccountConnected() && (
                    getprojectDetailsState(usePlugin.state).projectID ?
                        <ProjectDetails title={getprojectDetailsState(usePlugin.state).projectName}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
