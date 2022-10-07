/**
 * Keep all plugin related constants here
*/
import {
    AzureDevops,
    HeaderCSRFToken,
    MMCSRF,
    pluginId,
    RightSidebarHeader,
    eventTypeMap,
    serviceTypeMap,
    defaultPage,
    defaultPerPageLimit,
    SubscriptionFilterCreatedBy,
} from './common';
import {SVGIcons} from './icons';
import {linkProjectModal, createTaskModal, subscriptionModal, subscriptionFilterOptions} from './form';
import {pluginApiServiceConfigs} from './apiService';
import {error} from './messages';

export default {
    common: {
        pluginId,
        MMCSRF,
        HeaderCSRFToken,
        AzureDevops,
        RightSidebarHeader,
        eventTypeMap,
        serviceTypeMap,
        defaultPage,
        defaultPerPageLimit,
        SubscriptionFilterCreatedBy,
    },
    form: {
        linkProjectModal,
        createTaskModal,
        subscriptionModal,
        subscriptionFilterOptions,
    },
    messages: {
        error,
    },
    pluginApiServiceConfigs,
    SVGIcons,
};
