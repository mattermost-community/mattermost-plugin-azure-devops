import {Store, Action} from 'redux';
import {GlobalState} from 'mattermost-redux/types/store';

import {toggleIsConnected, toggleIsSubscriptionDeleted} from 'reducers/websocketEvent';

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

// TODO: types are already added in a PR so, update that here when synced with that
export function handleSubscriptionDeleted(store: Store<GlobalState, Action<Record<string, unknown>>>) {
    return (_: any) => {
        store.dispatch(toggleIsSubscriptionDeleted(true) as Action);
    };
}
