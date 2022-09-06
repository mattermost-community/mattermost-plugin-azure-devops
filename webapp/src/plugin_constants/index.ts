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
} from './common';
import {subscriptionModal, createTaskModal} from './form';
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
    },
    form: {
        subscriptionModal,
        createTaskModal,
    },
    messages: {
        error,
    },
    pluginApiServiceConfigs,
};
