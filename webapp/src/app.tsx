import React, {useEffect} from 'react';

import Rhs from 'containers/Rhs';
import plugin_constants from 'plugin_constants';
import usePluginApi from 'hooks/usePluginApi';

// Global styles
import 'styles/main.scss';
import {getLinkModalState, getTaskModalState} from 'selectors';

/**
 * Mattermost plugin allows registering only one component in RHS
 * So, we would be grouping all the different components inside "Rhs" component to generate one final component for registration
 */
const App = (): JSX.Element => {
    const usePlugin = usePluginApi();

    useEffect(() => {
        if (getLinkModalState(usePlugin.state).visibility || getTaskModalState(usePlugin.state).visibility) {
            usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        }
    }, [getLinkModalState(usePlugin.state).visibility, getTaskModalState(usePlugin.state).visibility]);

    return (
        <></>
    );
};

export default App;
