import React, {useCallback, useEffect, useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type ReposFilterProps = {
    handleSelectRepo: (repo: string, repoName?: string) => void
    selectedRepo: string
    selectedTargetBranch: string
    handleSelectTargetBranch: (branch: string) => void
} & ReposSubscriptionFiltersRequest

const ReposFilter = ({handleSelectRepo, project, organization, selectedRepo, selectedTargetBranch, handleSelectTargetBranch}: ReposFilterProps) => {
    const {
        getApiState,
        makeApiRequestWithCompletionStatus,
    } = usePluginApi();

    const reposSubscriptionFiltersRequest = useMemo<ReposSubscriptionFiltersRequest>(() => ({
        organization,
        project,
    }), [organization, project]);

    const reposSubscriptionTargetBranchFiltersRequest = useMemo<ReposSubscriptionFiltersRequest>(() => ({
        organization,
        project,
        repository: selectedRepo,
    }), [organization, project, selectedRepo]);

    useEffect(() => {
        makeApiRequestWithCompletionStatus(
            pluginConstants.pluginApiServiceConfigs.getRepositories.apiServiceName,
            reposSubscriptionFiltersRequest,
        );
    }, [reposSubscriptionFiltersRequest]);

    useEffect(() => {
        if (selectedRepo !== filterLabelValuePairAll.value) {
            makeApiRequestWithCompletionStatus(
                pluginConstants.pluginApiServiceConfigs.getRepositoryBranches.apiServiceName,
                reposSubscriptionTargetBranchFiltersRequest,
            );
        }
    }, [reposSubscriptionTargetBranchFiltersRequest.repository]);

    const {data: repositoriesDataFromApi, isLoading: isGetRepositoriesLoading, isError: isGetRepositoriesError} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getRepositories.apiServiceName,
        reposSubscriptionFiltersRequest as APIRequestPayload,
    );
    const repositoriesData = repositoriesDataFromApi as ReposSubscriptionFiltersResponse[] || [];

    const {data: repositoryBranchesDataFromApi, isLoading: isGetRepositoryBranchesLoading, isError: isGetRepositoryBranchesError} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getRepositoryBranches.apiServiceName,
        reposSubscriptionTargetBranchFiltersRequest as APIRequestPayload,
    );
    const repositoryBranchesData = repositoryBranchesDataFromApi as ReposSubscriptionTargetBranchFilterResponse[] || [];

    const getTargetBranchOptions = () => ([{...filterLabelValuePairAll}, ...formLabelValuePairs('name', 'name', repositoryBranchesData)]);

    return (
        <>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Repository'
                    value={selectedRepo}
                    onChange={handleSelectRepo}
                    options={[{...filterLabelValuePairAll}, ...formLabelValuePairs('name', 'id', repositoriesData)]}
                    error={isGetRepositoriesError}
                    loadingOptions={isGetRepositoriesLoading}
                    disabled={isGetRepositoriesLoading}
                />
            </div>
            <Dropdown
                placeholder='Target Branch'
                value={selectedTargetBranch}
                onChange={handleSelectTargetBranch}
                options={getTargetBranchOptions()}
                error={isGetRepositoryBranchesError}
                loadingOptions={isGetRepositoryBranchesLoading}
                disabled={selectedRepo === filterLabelValuePairAll.value || isGetRepositoryBranchesLoading}
            />
        </>
    );
};

export default ReposFilter;
