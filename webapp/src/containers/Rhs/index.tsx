import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getProjectDetailsState} from 'selectors';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {isUserAccountConnected, state} = usePluginApi();

    return (
        <div className='overflow-auto height-rhs bg-sidebar padding-25'>
            {!isUserAccountConnected() && <AccountNotLinked/>}
            {
                isUserAccountConnected() && (
                    getProjectDetailsState(state).projectID ?
                        <ProjectDetails {...getProjectDetailsState(state)}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
