/**
 * Keep all plugin related constants here
*/
import {
    AzureDevops,
    HeaderCSRFToken,
    MMCSRF,
    pluginId,
    RightSidebarHeader,
    boardsEventTypeMap,
    channelType,
} from './common';
import {SVGIcons} from './icons';
import {linkProjectModal, createTaskModal, subscriptionModal} from './form';
import {pluginApiServiceConfigs} from './apiService';
import {error} from './messages';

export default {
    common: {
        pluginId,
        MMCSRF,
        HeaderCSRFToken,
        AzureDevops,
        RightSidebarHeader,
        boardsEventTypeMap,
        channelType,
    },
    form: {
        linkProjectModal,
        createTaskModal,
        subscriptionModal,
    },
    messages: {
        error,
    },
    pluginApiServiceConfigs,
    SVGIcons,
};
