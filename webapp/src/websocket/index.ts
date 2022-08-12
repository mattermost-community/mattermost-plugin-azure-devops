import {Store, Action} from 'redux';
import {GlobalState} from 'mattermost-redux/types/store';

import {toggleIsConnected} from 'reducers/userConnected';

export function handleConnect(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: any) => {
        store.dispatch(toggleIsConnected(true) as Action);
    };
}

export function handleDisconnect(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: any) => {
        store.dispatch(toggleIsConnected(false) as Action);
    };
}
