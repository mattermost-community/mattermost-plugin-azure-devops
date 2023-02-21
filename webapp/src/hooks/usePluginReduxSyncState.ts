import {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import pluginConstants from 'pluginConstants';

import {resetGlobalModalState} from 'reducers/globalModal';
import {toggleIsLinkedProjectListChanged, toggleShowLinkModal} from 'reducers/linkModal';
import {toggleShowSubscribeModal} from 'reducers/subscribeModal';
import {toggleShowTaskModal} from 'reducers/taskModal';
import {getGlobalModalState, getLinkModalState, getSubscribeModalState, getCreateTaskModalState, getRhsState, getWebsocketEventState} from 'selectors';

import usePluginApi from 'hooks/usePluginApi';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

/**
 * This hook intercepts the redux changes for the plugin to perform common actions like
 * checking if the user's account is connected
 * closing any open modal if the user is logged out from another client tab
 * preventing opening a modal if the user's account is not connected
 * fetching linked project list on the opening of RHS, on linking of a new project, etc
 */
function usePluginReduxSyncState() {
    const {state, makeApiRequestWithCompletionStatus} = usePluginApi();
    const dispatch = useDispatch();

    const {isConnected} = getWebsocketEventState(state);
    const {modalId, commandArgs} = getGlobalModalState(state);
    const {isSidebarOpen} = getRhsState(state);
    const {visibility: linkProjectModalVisibility, isLinked} = getLinkModalState(state);
    const {visibility: createTaskModalVisibility} = getCreateTaskModalState(state);
    const {visibility: subscribeModalVisibility} = getSubscribeModalState(state);

    // Check if user is connected
    useEffect(() => {
        if (!isConnected) {
            makeApiRequestWithCompletionStatus(pluginConstants.pluginApiServiceConfigs.getUserDetails.apiServiceName);
        }
    }, [isSidebarOpen, modalId]);

    useApiRequestCompletionState({
        serviceName: pluginConstants.pluginApiServiceConfigs.getUserDetails.apiServiceName,
        handleError: () => dispatch(resetGlobalModalState()),
    });

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
        }

        if (!isConnected) {
            dispatch(toggleShowLinkModal({isVisible: false, commandArgs}));
            dispatch(toggleShowSubscribeModal({isVisible: false, commandArgs}));
            dispatch(toggleShowTaskModal({isVisible: false, commandArgs}));
        }
    }, [modalId, isConnected]);

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
                pluginConstants.pluginApiServiceConfigs.getAllLinkedProjectsList.apiServiceName,
            );
        }
    }, [
        isLinked,
        isSidebarOpen,
        isConnected,
        createTaskModalVisibility,
        subscribeModalVisibility,
    ]);
}

export default usePluginReduxSyncState;
