import React, {useCallback, useEffect, useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type BoardsFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedAreaPath: string
    handleSelectAreaPath: (value: string, name?: string) => void
    setIsFiltersError: (value: boolean) => void
}

const BoardsFilter = ({
    organization,
    projectId,
    eventType,
    selectedAreaPath,
    handleSelectAreaPath,
    setIsFiltersError,
}: BoardsFilterProps) => {
    const {subscriptionFiltersNameForBoards, subscriptionFiltersForBoards} = pluginConstants.form;

    const {
        getApiState,
        makeApiRequestWithCompletionStatus,
    } = usePluginApi();

    const getSubscriptionFiltersRequest = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForBoards,
        eventType,
    }), [organization, projectId, eventType, subscriptionFiltersForBoards]);

    useEffect(() => {
        if (eventType) {
            makeApiRequestWithCompletionStatus(
                pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
                getSubscriptionFiltersRequest,
            );
        }
    }, [getSubscriptionFiltersRequest]);

    const {data, isLoading, isError, isSuccess} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getSubscriptionFilters.apiServiceName,
        getSubscriptionFiltersRequest as APIRequestPayload,
    );
    const filtersData = data as GetSubscriptionFiltersResponse || [];

    useEffect(() => {
        if (isError && !isSuccess) {
            setIsFiltersError(true);
        } else {
            setIsFiltersError(false);
        }
    }, [isLoading, isError, isSuccess]);

    const getAreaPathOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForBoards.areaPath], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);

    return (
        <Dropdown
            placeholder='Area Path'
            value={selectedAreaPath}
            onChange={handleSelectAreaPath}
            options={getAreaPathOptions()}
            error={isError}
            loadingOptions={isLoading}
            disabled={isLoading}
        />
    );
};

export default BoardsFilter;
