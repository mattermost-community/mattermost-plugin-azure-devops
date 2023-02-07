import {useCallback, useEffect} from 'react';

import {filterLabelValuePairAll} from 'pluginConstants/common';
import {formLabelValuePairs} from 'utils';

import pluginConstants from 'pluginConstants';

import usePluginApi from './usePluginApi';

type UseLoadFiltersParams = {
    isModalOpen: boolean
    setIsFiltersError: (value: boolean) => void
    getSubscriptionFiltersRequestParams: GetSubscriptionFiltersRequest
}

function useLoadFilters({isModalOpen, setIsFiltersError, getSubscriptionFiltersRequestParams}: UseLoadFiltersParams) {
    const {organization, projectId, eventType} = getSubscriptionFiltersRequestParams;
    const {getApiState, makeApiRequestWithCompletionStatus} = usePluginApi();

    useEffect(() => {
        if (isModalOpen && organization && projectId && eventType) {
            makeApiRequestWithCompletionStatus(
                pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
                getSubscriptionFiltersRequestParams,
            );
        }
    }, [getSubscriptionFiltersRequestParams]);

    const {data, isLoading, isError, isSuccess} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
        getSubscriptionFiltersRequestParams as APIRequestPayload,
    );

    const filtersData = data as GetSubscriptionFiltersResponse || [];

    const getFilterOptions = useCallback((filterList: Record<string, string>[]) => (
        !isLoading && isSuccess ?
            ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filterList, ['[Any]'])]) :
            [pluginConstants.common.filterLabelValuePairAll]
    ), [isLoading, isSuccess]);

    useEffect(() => {
        if (isError && !isSuccess) {
            setIsFiltersError(true);
        } else {
            setIsFiltersError(false);
        }
    }, [isError, isSuccess]);

    return {filtersData, isLoading, isError, isSuccess, getFilterOptions};
}

export default useLoadFilters;
