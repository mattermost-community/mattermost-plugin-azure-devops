import React, {useCallback, useEffect, useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type ReposFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedRepo: string
    handleSelectRepo: (value: string, name?: string) => void
    selectedTargetBranch: string
    handleSelectTargetBranch: (value: string, name?: string) => void
    selectedPullRequestCreatedBy: string
    handleSelectPullRequestCreatedBy: (value: string, name?: string) => void
    selectedPullRequestReviewersContains: string
    handlePullRequestReviewersContains: (value: string, name?: string) => void
    selectedPushedBy: string
    handleSelectPushedBy: (value: string, name?: string) => void
    selectedMergeResult: string
    handleSelectMergeResult: (value: string, name?: string) => void
    selectedNotificationType: string
    handleSelectNotificationType: (value: string, name?: string) => void
    setIsFiltersError: (value: boolean) => void
}

// TODO: Refactor the props to minimal
const ReposFilter = ({
    organization,
    projectId,
    eventType,
    selectedRepo,
    handleSelectRepo,
    selectedTargetBranch,
    handleSelectTargetBranch,
    selectedPullRequestCreatedBy,
    handleSelectPullRequestCreatedBy,
    selectedPullRequestReviewersContains,
    handlePullRequestReviewersContains,
    selectedPushedBy,
    handleSelectPushedBy,
    selectedMergeResult,
    handleSelectMergeResult,
    selectedNotificationType,
    handleSelectNotificationType,
    setIsFiltersError,
}: ReposFilterProps) => {
    const {mergeResultOptons, pullRequestChangeOptons, subscriptionFiltersForRepos, subscriptionFiltersNameForRepos} = pluginConstants.form;

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

    useEffect(() => {
        if (isError && !isSuccess) {
            setIsFiltersError(true);
        } else {
            setIsFiltersError(false);
        }
    }, [isLoading, isError, isSuccess]);

    const getRepositoryOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.repository], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getTargetBranchOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.branch], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getPullRequestCreatedByOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.pullrequestCreatedBy], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getPullRequestReviewersContainsOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.pullrequestReviewersContains], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getPullRequestPushedByOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForRepos.pushedBy], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);

    return (
        <>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Repository'
                    value={selectedRepo}
                    onChange={handleSelectRepo}
                    options={getRepositoryOptions()}
                    error={isError}
                    loadingOptions={isLoading}
                    disabled={isLoading}
                />
            </div>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Target Branch'
                    value={selectedTargetBranch}
                    onChange={handleSelectTargetBranch}
                    options={getTargetBranchOptions()}
                    error={isError}
                    loadingOptions={isLoading}
                    disabled={selectedRepo === filterLabelValuePairAll.value || isLoading}
                />
            </div>
            {
                eventType === pluginConstants.common.eventTypeReposKeys.merged && (
                    <div className='margin-bottom-10'>
                        <Dropdown
                            placeholder='Merge Result'
                            value={selectedMergeResult}
                            onChange={handleSelectMergeResult}
                            options={mergeResultOptons}
                            error={isError}
                            loadingOptions={isLoading}
                            disabled={isLoading}
                        />
                    </div>
                )
            }
            {
                eventType === pluginConstants.common.eventTypeReposKeys.updated && (
                    <div className='margin-bottom-10'>
                        <Dropdown
                            placeholder='Change'
                            value={selectedNotificationType}
                            onChange={handleSelectNotificationType}
                            options={pullRequestChangeOptons}
                            error={isError}
                            loadingOptions={isLoading}
                            disabled={isLoading}
                        />
                    </div>
                )
            }
            {
                eventType !== pluginConstants.common.eventTypeReposKeys.commented &&
                    eventType !== pluginConstants.common.eventTypeReposKeys.codePushed && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Requested by a member of group'
                                value={selectedPullRequestCreatedBy}
                                onChange={handleSelectPullRequestCreatedBy}
                                options={getPullRequestCreatedByOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <Dropdown
                            placeholder='Reviewer includes group'
                            value={selectedPullRequestReviewersContains}
                            onChange={handlePullRequestReviewersContains}
                            options={getPullRequestReviewersContainsOptions()}
                            error={isError}
                            loadingOptions={isLoading}
                            disabled={isLoading}
                        />
                    </>
                )
            }
            {
                eventType === pluginConstants.common.eventTypeReposKeys.codePushed && (
                    <Dropdown
                        placeholder='Pushed by a member of group'
                        value={selectedPushedBy}
                        onChange={handleSelectPushedBy}
                        options={getPullRequestPushedByOptions()}
                        error={isError}
                        loadingOptions={isLoading}
                        disabled={isLoading}
                    />
                )
            }
        </>
    );
};

export default ReposFilter;
