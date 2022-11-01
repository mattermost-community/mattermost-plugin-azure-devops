import React, {useCallback, useEffect, useMemo, useState} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type ReposFilterProps = {
    organization: string;
    projectId: string;
    eventType: string;
    selectedRepo: string
    handleSelectRepo: (repo: string, repoName?: string) => void
    selectedTargetBranch: string
    handleSelectTargetBranch: (branch: string) => void
    selectedPullRequestCreatedBy: string
    handleSelectPullRequestCreatedBy: (pullrequestCreatedBy: string) => void
    selectedPullRequestReviewersContains: string
    handlePullRequestReviewersContains: (pullrequestReviewersContains: string) => void
}

const subscriptionFiltersNameForRepos = {
    repository: 'repository',
    branch: 'branch',
    pullrequestCreatedBy: 'pullrequestCreatedBy',
    pullrequestReviewersContains: 'pullrequestReviewersContains',
};
const subscriptionFiltersForRepos = [
    subscriptionFiltersNameForRepos.repository,
    subscriptionFiltersNameForRepos.branch,
    subscriptionFiltersNameForRepos.pullrequestCreatedBy,
    subscriptionFiltersNameForRepos.pullrequestReviewersContains,
];

const ReposFilter = ({organization, projectId, eventType, selectedRepo, handleSelectRepo, selectedTargetBranch, handleSelectTargetBranch, selectedPullRequestCreatedBy, handleSelectPullRequestCreatedBy, selectedPullRequestReviewersContains, handlePullRequestReviewersContains}: ReposFilterProps) => {
    const {
        getApiState,
        makeApiRequestWithCompletionStatus,
    } = usePluginApi();

    const getSubscriptionFiltersRequest = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForRepos,
        eventType,
        repositoryId: selectedRepo,
    }), [organization, projectId, eventType, subscriptionFiltersForRepos, selectedRepo]);

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

    const getRepositoryOptions = () => isSuccess && ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.repository], ['[Any]'])]);
    const getTargetBranchOptions = () => isSuccess && ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.branch], ['[Any]'])]);
    const getPullrequestCreatedByOptions = () => isSuccess && ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.pullrequestCreatedBy], ['[Any]'])]);
    const getPullrequestReviewersContainsOptions = () => isSuccess && ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.pullrequestReviewersContains], ['[Any]'])]);

    return (
        <>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Repository'
                    value={selectedRepo}
                    onChange={handleSelectRepo}
                    options={getRepositoryOptions() || [pluginConstants.common.filterLabelValuePairAll]}
                    error={isError}
                    loadingOptions={isLoading}
                    disabled={!eventType || isLoading}
                />
            </div>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Target Branch'
                    value={selectedTargetBranch}
                    onChange={handleSelectTargetBranch}
                    options={getTargetBranchOptions() || [pluginConstants.common.filterLabelValuePairAll]}
                    error={isError}
                    loadingOptions={isLoading}
                    disabled={selectedRepo === filterLabelValuePairAll.value || isLoading}
                />
            </div>
            {
                eventType !== pluginConstants.common.eventTypeReposKeys.commented &&
                eventType !== pluginConstants.common.eventTypeReposKeys.codePushed && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Requested by a member of group'
                                value={selectedPullRequestCreatedBy}
                                onChange={handleSelectPullRequestCreatedBy}
                                options={getPullrequestCreatedByOptions() || [pluginConstants.common.filterLabelValuePairAll]}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={!eventType || isLoading}
                            />
                        </div>
                        <Dropdown
                            placeholder='Reviewer includes group'
                            value={selectedPullRequestReviewersContains}
                            onChange={handlePullRequestReviewersContains}
                            options={getPullrequestReviewersContainsOptions() || [pluginConstants.common.filterLabelValuePairAll]}
                            error={isError}
                            loadingOptions={isLoading}
                            disabled={!eventType || isLoading}
                        />
                    </>
                )
            }
        </>
    );
};

export default ReposFilter;
