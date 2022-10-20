import React, {useEffect, useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type ReposFilterProps = {
    handleSelectRepo: (repo: string) => void
    selectedRepo: string
} & ReposSubscriptionFiltersRequest

const ReposFilter = ({handleSelectRepo, project, organization, selectedRepo}: ReposFilterProps) => {
    const {
        getApiState,
        makeApiRequestWithCompletionStatus,
    } = usePluginApi();

    const reposSubscriptionFiltersRequest = useMemo<ReposSubscriptionFiltersRequest>(() => ({
        organization,
        project,
    }), [organization, project]);

    useEffect(() => {
        makeApiRequestWithCompletionStatus(
            pluginConstants.pluginApiServiceConfigs.getRepositories.apiServiceName,
            reposSubscriptionFiltersRequest,
        );
    }, [reposSubscriptionFiltersRequest]);

    const {data, isLoading: isGetRepositoriesLoading, isError: isGetRepositoriesError} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getRepositories.apiServiceName,
        reposSubscriptionFiltersRequest as APIRequestPayload,
    );
    const repositoriesData = data as ReposSubscriptionFiltersResponse[] || [];

    return (
        <Dropdown
            placeholder='Repository'
            value={selectedRepo}
            onChange={handleSelectRepo}
            options={[filterLabelValuePairAll, ...formLabelValuePairs('name', 'id', repositoriesData)]}
            error={isGetRepositoriesError}
            loadingOptions={isGetRepositoriesLoading}
        />
    );
};

export default ReposFilter;
