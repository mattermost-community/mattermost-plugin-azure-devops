import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getProjectDetailsState, getRhsState} from 'selectors';

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
                    getProjectDetailsState(usePlugin.state).projectID ?
                        <ProjectDetails {...getProjectDetailsState(usePlugin.state)}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
