import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import reducer from 'reducers';

import {handleConnect, handleDisconnect, handleSubscriptionDeleted} from 'websocket';

import {ChannelHeaderBtn} from 'components/buttons/action_buttons';

import Constants from 'plugin_constants';

import Hooks from 'hooks';

import Rhs from 'containers/Rhs';
import LinkModal from 'containers/modals/LinkModal';
import TaskModal from 'containers/modals/TaskModal';
import SubscribeModal from 'containers/modals/SubscribeModal';

import App from './app';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';
import manifest from './manifest';

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(App);
        registry.registerRootComponent(TaskModal);
        registry.registerRootComponent(LinkModal);
        registry.registerRootComponent(SubscribeModal);
        registry.registerWebSocketEventHandler(`custom_${Constants.common.pluginId}_connect`, handleConnect(store));
        registry.registerWebSocketEventHandler(`custom_${Constants.common.pluginId}_disconnect`, handleDisconnect(store));
        registry.registerWebSocketEventHandler(`custom_${Constants.common.pluginId}_subscription_deleted`, handleSubscriptionDeleted(store));
        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(Rhs, Constants.common.RightSidebarHeader);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
        registry.registerChannelHeaderButtonAction(<ChannelHeaderBtn/>, () => store.dispatch(showRHSPlugin), null, Constants.common.AzureDevops);
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
