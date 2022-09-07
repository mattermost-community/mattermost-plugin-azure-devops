import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import plugin_constants from 'plugin_constants';

import {resetGlobalModalState} from 'reducers/globalModal';
import {toggleIsLinkedProjectListChanged, toggleShowLinkModal} from 'reducers/linkModal';
import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowTaskModal} from 'reducers/taskModal';
import {getGlobalModalState, getLinkModalState, getSubscribeModalState, getCreateTaskModalState, getRhsState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';

// Global styles
import 'styles/main.scss';

/**
 * This is a central component for adding account connection validation on all the modals registered in the root component
 */
const App = (): JSX.Element => {
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    // Check if user is connected on page reload
    useEffect(() => {
        usePlugin.makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
    }, []);

    /**
     * When a command is issued on the Mattermost to open any modal
     * then here we first check if the user's account is connected or not
     * if the account is connected we dispatch the action to open the required modal
     * otherwise we reset the action and don't open any modal
     */
    useEffect(() => {
        const {modalId, commandArgs} = getGlobalModalState(usePlugin.state);

        if (usePlugin.isUserAccountConnected() && modalId) {
            switch (modalId) {
            case 'linkProject':
                dispatch(toggleShowLinkModal({isVisible: true, commandArgs}));
                break;
            case 'subscribeProject':
                dispatch(toggleShowSubscribeModal({isVisible: true, commandArgs}));
                break;
            case 'createBoardTask':
                dispatch(toggleShowTaskModal({isVisible: true, commandArgs}));
            }
        } else {
            dispatch(resetGlobalModalState());
        }
    }, [getGlobalModalState(usePlugin.state).modalId]);

    useEffect(() => {
        dispatch(resetGlobalModalState());
    }, [
        getLinkModalState(usePlugin.state).visibility,
        getCreateTaskModalState(usePlugin.state).visibility,
        getSubscribeModalState(usePlugin.state).visibility,
    ]);

    // Fetch the list of linked projects
    useEffect(() => {
        if (usePlugin.isUserAccountConnected()) {
            if (getLinkModalState(usePlugin.state).isLinked) {
                dispatch(toggleIsLinkedProjectListChanged(false));
            }

            usePlugin.makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
            );
        }
    }, [
        usePlugin.isUserAccountConnected(),
        getCreateTaskModalState(usePlugin.state).visibility,
        getSubscribeModalState(usePlugin.state).visibility,
        getRhsState(usePlugin.state).isSidebarOpen,
        getLinkModalState(usePlugin.state).isLinked,
    ]);

    return <></>;
};

export default App;
