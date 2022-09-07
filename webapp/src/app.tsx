import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import plugin_constants from 'plugin_constants';

import {resetGlobalModalState} from 'reducers/globalModal';
import {toggleIsLinkedProjectListChanged, toggleShowLinkModal} from 'reducers/linkModal';
import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowTaskModal} from 'reducers/taskModal';
import {getGlobalModalState, getLinkModalState, getSubscribeModalState, getCreateTaskModalState, getRhsState, getWebsocketEventState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';

// Global styles
import 'styles/main.scss';

/**
 * This is a central component for adding account connection validation on all the modals registered in the root component
 */
const App = (): JSX.Element => {
    const {state, makeApiRequest, makeApiRequestWithCompletionStatus} = usePluginApi();
    const dispatch = useDispatch();

    const {isConnected} = getWebsocketEventState(state);
    const {modalId, commandArgs} = getGlobalModalState(state);
    const {isSidebarOpen} = getRhsState(state);
    const {visibility: linkProjectModalVisibility, isLinked} = getLinkModalState(state);
    const {visibility: createTaskModalVisibility} = getCreateTaskModalState(state);
    const {visibility: subscribeModalVisibility} = getSubscribeModalState(state);

    // Check if user is connected on page reload
    useEffect(() => {
        makeApiRequest(plugin_constants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
    }, []);

    /**
     * When a command is issued on the Mattermost to open any modal
     * then here we first check if the user's account is connected or not
     * if the account is connected we dispatch the action to open the required modal
     * otherwise we reset the action and don't open any modal
     */
    useEffect(() => {
        if (isConnected && modalId) {
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
    }, [modalId]);

    useEffect(() => {
        dispatch(resetGlobalModalState());
    }, [
        linkProjectModalVisibility,
        createTaskModalVisibility,
        subscribeModalVisibility,
    ]);

    // Fetch the list of linked projects
    useEffect(() => {
        if (isConnected) {
            if (isLinked) {
                dispatch(toggleIsLinkedProjectListChanged(false));
            }

            makeApiRequestWithCompletionStatus(
                plugin_constants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
            );
        }
    }, [
        isLinked,
        isSidebarOpen,
        isConnected,
        createTaskModalVisibility,
        subscribeModalVisibility,
    ]);

    return <></>;
};

export default App;
