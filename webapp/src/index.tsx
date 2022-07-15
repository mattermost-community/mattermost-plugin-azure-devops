import React from 'react';
import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import reducer from 'reducers';

import Rhs from 'containers/Rhs';
import {ChannelHeaderBtn, ChannelIntroBtn} from 'containers/action_buttons';

import Constants from 'plugin_constants';

import manifest from './manifest';

import App from './app';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        registry.registerReducer(reducer);
        registry.registerRootComponent(App);
        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(Rhs, Constants.RightSidebarHeader);

        registry.registerChannelHeaderButtonAction(<ChannelHeaderBtn/>, () => store.dispatch(showRHSPlugin), null, Constants.AzureDevops);
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
