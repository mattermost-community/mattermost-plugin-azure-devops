import React, {useEffect} from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState, getRhsState, getUserConnectionState} from 'selectors';

import LinearLoader from 'components/loader/linear';

import plugin_constants from 'plugin_constants';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const usePlugin = usePluginApi();

    // Fetch the connected account details when RHS is opened
    useEffect(() => {
        if (getRhsState(usePlugin.state).isSidebarOpen) {
            usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        }
    }, []);

    if (!getRhsState(usePlugin.state).isSidebarOpen) {
        return <></>;
    }

    return (
        <div className='overflow-auto height-rhs bg-sidebar padding-25'>
            {
                usePlugin.getApiState(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName).isLoading &&
                <LinearLoader/>
            }
            {
                !usePlugin.getUserAccountConnectionState().isLoading &&
                usePlugin.getUserAccountConnectionState().isError &&
                <AccountNotLinked/>
            }
            {
                !usePlugin.getUserAccountConnectionState().isLoading &&
                usePlugin.getUserAccountConnectionState().isSuccess &&
                usePlugin.getUserAccountConnectionState().data?.MattermostUserID && (
                    getprojectDetailsState(usePlugin.state).projectID ?
                        <ProjectDetails title={getprojectDetailsState(usePlugin.state).projectName}/> :
                        <ProjectList/>
                )
            }
        </div>
    );
};

export default Rhs;
