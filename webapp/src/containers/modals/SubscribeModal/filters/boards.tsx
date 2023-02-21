import React, {useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import Dropdown from 'components/dropdown';

import useLoadFilters from 'hooks/useLoadFilters';

type BoardsFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedAreaPath: string
    handleSetFilter: HandleSetSubscriptionFilter
    setIsFiltersError: (value: boolean) => void
    isModalOpen: boolean
}

const BoardsFilter = ({
    organization,
    projectId,
    eventType,
    selectedAreaPath,
    handleSetFilter,
    setIsFiltersError,
    isModalOpen,
}: BoardsFilterProps) => {
    const {subscriptionFiltersNameForBoards, subscriptionFiltersForBoards} = pluginConstants.form;

    const getSubscriptionFiltersRequestParams = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForBoards,
        eventType,
    }), [organization, projectId, eventType, subscriptionFiltersForBoards]);

    const {filtersData, isError, isLoading, getFilterOptions} = useLoadFilters({isModalOpen, getSubscriptionFiltersRequestParams, setIsFiltersError});

    return (
        <Dropdown
            placeholder='Area Path'
            value={selectedAreaPath}
            onChange={(newValue) => handleSetFilter('areaPath', newValue)}
            options={getFilterOptions(filtersData[subscriptionFiltersNameForBoards.areaPath])}
            error={isError}
            loadingOptions={isLoading}
            disabled={isLoading}
        />
    );
};

export default BoardsFilter;
