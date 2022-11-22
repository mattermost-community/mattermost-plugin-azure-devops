import React, {useEffect} from 'react';

import usePluginApi from 'hooks/usePluginApi';

import {getProjectDetailsState, getWebsocketEventState} from 'selectors';

import pluginConstants from 'pluginConstants';

import LinearLoader from 'components/loader/linear';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {state, makeApiRequest, getApiState} = usePluginApi();
    const {isConnected} = getWebsocketEventState(state);
    const {isLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.getUserDetails.apiServiceName);

    // Check if user is connected on page reload
    useEffect(() => {
        makeApiRequest(pluginConstants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
    }, []);

    return (
        <div
            id='scrollableArea'
            className='overflow-auto height-rhs position-relative padding-16'
        >
            {isLoading && <LinearLoader/>}
            {!isLoading && !isConnected && <AccountNotLinked/>}
            {
                !isLoading && isConnected && (
                    getProjectDetailsState(state).projectID ?
                        <ProjectDetails {...getProjectDetailsState(state)}/> :
                        <ProjectList/>)
            }
        </div>
    );
};

export default Rhs;
