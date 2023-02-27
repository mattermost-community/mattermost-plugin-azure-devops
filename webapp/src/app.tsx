import React from 'react';

// Hook
import usePluginReduxSyncState from 'hooks/usePluginReduxSyncState';

// Global styles
import 'styles/main.scss';

/**
 * This is a central component for intercepting the redux changes for the plugin using the hook "usePluginReduxSyncState"
 */
const App = (): JSX.Element => {
    usePluginReduxSyncState();

    return <></>;
};

export default App;
