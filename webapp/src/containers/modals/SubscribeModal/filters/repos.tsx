import React, {useCallback, useEffect, useMemo, useState} from 'react';

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

    const [targetBranch, setTargetBranch] = useState<string>(filterLabelValuePairAll.value);

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
        makeApiRequestWithCompletionStatus(
            pluginConstants.pluginApiServiceConfigs.getRepositoryBranches.apiServiceName,
            reposSubscriptionTargetBranchFiltersRequest,
        );
    }, [reposSubscriptionTargetBranchFiltersRequest]);

    const {data: repositoriesDateFromApi, isLoading: isGetRepositoriesLoading, isError: isGetRepositoriesError} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getRepositories.apiServiceName,
        reposSubscriptionFiltersRequest as APIRequestPayload,
    );
    const repositoriesData = repositoriesDateFromApi as ReposSubscriptionFiltersResponse[] || [];

    const {data: repositoryBranchesDataFromApi, isLoading: isGetRepositoryBranchesLoading, isError: isGetRepositoryBranchesError} = getApiState(
        pluginConstants.pluginApiServiceConfigs.getRepositoryBranches.apiServiceName,
        reposSubscriptionTargetBranchFiltersRequest as APIRequestPayload,
    );
    const repositoryBranchesData = repositoryBranchesDataFromApi as ReposSubscriptionTargetBranchFilterResponse[] || [];

    const getTargetBranchOptions = useCallback(() => ([{...filterLabelValuePairAll}, ...formLabelValuePairs('name', 'objectId', repositoryBranchesData)]), [selectedRepo]);

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
                value={targetBranch}
                onChange={(newValue) => setTargetBranch(newValue)}
                options={getTargetBranchOptions()}
                error={isGetRepositoryBranchesError}
                loadingOptions={isGetRepositoryBranchesLoading}
                disabled={selectedRepo === filterLabelValuePairAll.value || isGetRepositoryBranchesLoading}
            />
        </>
    );
};

export default ReposFilter;
