import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import reducer from 'reducers';

import {handleConnect, handleDisconnect} from 'websocket';

import {ChannelHeaderBtn} from 'containers/action_buttons';

import Constants from 'plugin_constants';

import Hooks from 'hooks';

import Rhs from 'containers/Rhs';
import LinkModal from 'containers/LinkModal';
import TaskModal from 'containers/TaskModal';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';
import App from './app';
import manifest from './manifest';

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(App);
        registry.registerRootComponent(TaskModal);
        registry.registerRootComponent(LinkModal);
        registry.registerWebSocketEventHandler(`custom_${Constants.pluginId}_connect`, handleConnect(store));
        registry.registerWebSocketEventHandler(`custom_${Constants.pluginId}_disconnect`, handleDisconnect(store));
        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(Rhs, Constants.RightSidebarHeader);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
        registry.registerChannelHeaderButtonAction(<ChannelHeaderBtn/>, () => store.dispatch(showRHSPlugin), null, Constants.AzureDevops);
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
