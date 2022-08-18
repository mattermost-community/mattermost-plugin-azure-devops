import React, {useEffect} from 'react';

import usePluginApi from 'hooks/usePluginApi';
import {getprojectDetailsState, getRhsState} from 'selectors';

import LinearLoader from 'components/loader/linear';

import plugin_constants from 'plugin_constants';

import AccountNotLinked from './accountNotLinked';
import ProjectList from './projectList';
import ProjectDetails from './projectDetails';

const Rhs = (): JSX.Element => {
    const {getApiState, makeApiRequest, state} = usePluginApi();

    // Fetch the connected account details when RHS is opened
    useEffect(() => {
        if (getRhsState(state).isSidebarOpen) {
            makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        }
    }, []);

    const {isLoading, isError, isSuccess} = getApiState(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);

    if (!getRhsState(state).isSidebarOpen) {
        return <></>;
    }

    return (
        <div className='overflow-auto height-rhs bg-sidebar padding-25'>
            {
                isLoading &&
                <LinearLoader/>
            }
            {
                isError &&
                <AccountNotLinked/>
            }
            {
                !isLoading &&
                !isError &&
                isSuccess && (
                    getprojectDetailsState(state).projectID ?
                        <ProjectDetails title={getprojectDetailsState(state).projectName}/> :
                        <ProjectList/>
                )

            }
        </div>
    );
};

export default Rhs;
