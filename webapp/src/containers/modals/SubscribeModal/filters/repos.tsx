import React, {useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import useLoadFilters from 'hooks/useLoadFilters';

type ReposFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedRepo: string
    isModalOpen: boolean
    handleSetFilter: HandleSetSubscriptionFilter
    selectedTargetBranch: string
    selectedPullRequestCreatedBy: string
    selectedPullRequestReviewersContains: string
    selectedPushedBy: string
    selectedMergeResult: string
    selectedNotificationType: string
    setIsFiltersError: (value: boolean) => void
}

// TODO: Refactor the props to minimal
const ReposFilter = ({
    organization,
    projectId,
    eventType,
    selectedRepo,
    isModalOpen,
    handleSetFilter,
    selectedTargetBranch,
    selectedPullRequestCreatedBy,
    selectedPullRequestReviewersContains,
    selectedPushedBy,
    selectedMergeResult,
    selectedNotificationType,
    setIsFiltersError,
}: ReposFilterProps) => {
    const {mergeResultOptons, pullRequestChangeOptons, subscriptionFiltersForRepos, subscriptionFiltersNameForRepos} = pluginConstants.form;

    const getSubscriptionFiltersRequestParams = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForRepos,
        eventType,
        repositoryId: selectedRepo,
    }), [organization, projectId, eventType, subscriptionFiltersForRepos, selectedRepo]);

    const {filtersData, isError, isLoading, getFilterOptions} = useLoadFilters({isModalOpen, getSubscriptionFiltersRequestParams, setIsFiltersError});

    return (
        <>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Repository'
                    value={selectedRepo}
                    onChange={(newValue, label) => handleSetFilter('repository', newValue, 'repositoryName', label)}
                    options={getFilterOptions(filtersData[subscriptionFiltersNameForRepos.repository])}
                    error={isError}
                    loadingOptions={isLoading}
                    disabled={isLoading}
                />
            </div>
            <div className='margin-bottom-10'>
                <Dropdown
                    placeholder='Target Branch'
                    value={selectedTargetBranch}
                    onChange={(newValue) => handleSetFilter('targetBranch', newValue)}
                    options={getFilterOptions(filtersData[subscriptionFiltersNameForRepos.branch])}
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
                            onChange={(newValue, label) => handleSetFilter('mergeResult', newValue, 'mergeResultName', label)}
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
                            onChange={(newValue, label) => handleSetFilter('notificationType', newValue, 'notificationTypeName', label)}
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
                                onChange={(newValue, label) => handleSetFilter('pullRequestCreatedBy', newValue, 'pullRequestCreatedByName', label)}
                                options={getFilterOptions(filtersData[subscriptionFiltersNameForRepos.pullrequestCreatedBy])}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <Dropdown
                            placeholder='Reviewer includes group'
                            value={selectedPullRequestReviewersContains}
                            onChange={(newValue, label) => handleSetFilter('pullRequestReviewersContains', newValue, 'pullRequestReviewersContainsName', label)}
                            options={getFilterOptions(filtersData[subscriptionFiltersNameForRepos.pullrequestReviewersContains])}
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
                        onChange={(newValue, label) => handleSetFilter('pushedBy', newValue, 'pushedByName', label)}
                        options={getFilterOptions(filtersData[subscriptionFiltersNameForRepos.pushedBy])}
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
