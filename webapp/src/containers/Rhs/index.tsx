import React from 'react';

import usePluginApi from 'hooks/usePluginApi';

import {getProjectDetailsState, getWebsocketEventState} from 'selectors';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {state} = usePluginApi();
    const {isConnected} = getWebsocketEventState(state);

    return (
        <div
            id='scrollableArea'
            className='overflow-auto height-rhs position-relative padding-16'
        >
            {!isConnected && <AccountNotLinked/>}
            {
                isConnected && (
                    getProjectDetailsState(state).projectID ?
                        <ProjectDetails {...getProjectDetailsState(state)}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
