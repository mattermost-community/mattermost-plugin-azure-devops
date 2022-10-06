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
    defaultPage,
    defaultPerPageLimit,
    subscriptionFilters,
} from './common';
import {SVGIcons} from './icons';
import {linkProjectModal, createTaskModal, subscriptionModal, subscriptionFilterCreatedByOptions, subscriptionFilterServiceTypeOptions} from './form';
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
        defaultPage,
        defaultPerPageLimit,
        subscriptionFilters,
    },
    form: {
        linkProjectModal,
        createTaskModal,
        subscriptionModal,
        subscriptionFilterCreatedByOptions,
        subscriptionFilterServiceTypeOptions,
    },
    messages: {
        error,
    },
    pluginApiServiceConfigs,
    SVGIcons,
};
