import React from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState} from 'selectors';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {isUserAccountConnected, state} = usePluginApi();

    return (
        <div className='overflow-auto height-rhs bg-sidebar padding-25'>
            {
                !isUserAccountConnected() && <AccountNotLinked/>
            }
            {
                isUserAccountConnected() && (
                    getprojectDetailsState(state).projectID ?
                        <ProjectDetails title={getprojectDetailsState(state).projectName}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
