/**
 * Keep all plugin related constants here
*/
import {
    AzureDevops,
    HeaderCSRFToken,
    MMCSRF,
    pluginId,
    RightSidebarHeader,
    ToggleSwitchLabelPositioning,
    ToggleLabel,
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
        ToggleSwitchLabelPositioning,
        ToggleLabel,
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
