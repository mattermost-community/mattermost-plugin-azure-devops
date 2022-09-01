import React, {useEffect} from 'react';
import {useDispatch} from 'react-redux';

import usePluginApi from 'hooks/usePluginApi';

import {getGlobalModalState, getLinkModalState} from 'selectors';

import {toggleShowLinkModal} from 'reducers/linkModal';
import {resetGlobalModalState} from 'reducers/globalModal';

// Global styles
import 'styles/main.scss';

/**
 * This is a central component for adding account connection validation on all the modals registered in the root component
 */
const App = (): JSX.Element => {
    const usePlugin = usePluginApi();
    const dispatch = useDispatch();

    /**
     * When a command is issued on the Mattermost to open any modal
     * then here we first check if the user's account is connected or not
     * If the account is connected, we dispatch the action to open the required modal
     * otherwise we reset the action and don't open any modal
     */
    useEffect(() => {
        const {modalId, commandArgs} = getGlobalModalState(usePlugin.state);

        if (usePlugin.isUserAccountConnected() && modalId) {
            switch (modalId) {
            case 'linkProject':
                dispatch(toggleShowLinkModal({isVisible: true, commandArgs}));
                break;
            }
        } else {
            dispatch(resetGlobalModalState());
        }
    }, [getGlobalModalState(usePlugin.state).modalId]);

    useEffect(() => {
        dispatch(resetGlobalModalState());
    }, [getLinkModalState(usePlugin.state).visibility]);

    return <></>;
};

export default App;
