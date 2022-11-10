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
    subscriptionFilters,
    boards,
    repos,
    defaultSubscriptionFilters,
    filterLabelValuePairAll,
    eventTypeReposKeys,
} from './common';
import {SVGIcons} from './icons';
import {
    linkProjectModal,
    createTaskModal,
    subscriptionModal,
    subscriptionFilterCreatedByOptions,
    subscriptionFilterServiceTypeOptions,
    subscriptionFilterEventTypeBoardsOptions,
    subscriptionFilterEventTypeReposOptions,
    mergeResultOptons,
    pullRequestChangeOptons,
    subscriptionFiltersNameForRepos,
    subscriptionFiltersForRepos,
    subscriptionFiltersNameForBoards,
    subscriptionFiltersForBoards,
} from './form';
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
        subscriptionFilters,
        boards,
        repos,
        defaultSubscriptionFilters,
        filterLabelValuePairAll,
        eventTypeReposKeys,
    },
    form: {
        linkProjectModal,
        createTaskModal,
        subscriptionModal,
        subscriptionFilterCreatedByOptions,
        subscriptionFilterServiceTypeOptions,
        subscriptionFilterEventTypeBoardsOptions,
        subscriptionFilterEventTypeReposOptions,
        mergeResultOptons,
        pullRequestChangeOptons,
        subscriptionFiltersNameForRepos,
        subscriptionFiltersForRepos,
        subscriptionFiltersNameForBoards,
        subscriptionFiltersForBoards,
    },
    messages: {
        error,
    },
    pluginApiServiceConfigs,
    SVGIcons,
};
