import React from 'react';

import usePluginApi from 'hooks/usePluginApi';

import {getProjectDetailsState, getWebsocketEventState} from 'selectors';

import pluginConstants from 'pluginConstants';

import LinearLoader from 'components/loader/linear';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {state, getApiState} = usePluginApi();
    const {isConnected} = getWebsocketEventState(state);
    const {isLoading} = getApiState(pluginConstants.pluginApiServiceConfigs.getUserDetails.apiServiceName);

    if (isLoading) {
        return <LinearLoader/>;
    }

    return (
        <div
            id='scrollableArea'
        >
            {!isConnected && <AccountNotLinked/>}
            {isConnected && (
                getProjectDetailsState(state).projectID ?
                    <ProjectDetails {...getProjectDetailsState(state)}/> :
                    <ProjectList/>)}
        </div>
    );
};

export default Rhs;
