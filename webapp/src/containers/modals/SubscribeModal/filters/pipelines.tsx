import React, {useCallback, useEffect, useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import usePluginApi from 'hooks/usePluginApi';

import {formLabelValuePairs} from 'utils';

type PipelinesFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedBuildPipeline: string
    handleSelectBuildPipeline: (value: string, name?: string) => void
    setIsFiltersError: (value: boolean) => void
    selectedBuildStatus: string
    handleSelectBuildStatus: (value: string, name?: string) => void
    handleSelectReleasePipeline: (value: string, name?: string) => void
    selectedReleasePipeline: string
    handleSelectStageName: (value: string, name?: string) => void
    selectedStageName: string
    handleSelectApprovalType: (value: string, name?: string) => void
    selectedApprovalType: string
    handleSelectApprovalStatus: (value: string, name?: string) => void
    selectedApprovalStatus: string
    selectedReleaseStatus: string
    handleSelectReleaseStatus: (value: string, name?: string) => void
}

const PipelinesFilter = ({
    organization,
    projectId,
    eventType,
    selectedBuildPipeline,
    handleSelectBuildPipeline,
    setIsFiltersError,
    selectedBuildStatus,
    handleSelectBuildStatus,
    handleSelectReleasePipeline,
    selectedReleasePipeline,
    handleSelectStageName,
    selectedStageName,
    selectedApprovalType,
    handleSelectApprovalType,
    handleSelectApprovalStatus,
    selectedApprovalStatus,
    selectedReleaseStatus,
    handleSelectReleaseStatus,
}: PipelinesFilterProps) => {
    const {buildStatusOptions, releaseApprovalTypeOptions, releaseApprovalStatusOptions, releaseStatusOptions, subscriptionFiltersForPipelines, subscriptionFiltersNameForPipelines} = pluginConstants.form;

    const {
        getApiState,
        makeApiRequestWithCompletionStatus,
    } = usePluginApi();

    const getSubscriptionFiltersRequest = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForPipelines,
        eventType,
        releasePipelineId: selectedReleasePipeline,
    }), [organization, projectId, eventType, subscriptionFiltersForPipelines, selectedBuildPipeline, selectedReleasePipeline]);

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

    const getBuildPipelineOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.buildPipeline], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getReleasePipelineOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.releasePipelineName], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getStageNameOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.stageName], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);

    return (
        <>
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.buildCompleted && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Build Pipeline'
                                value={selectedBuildPipeline}
                                onChange={handleSelectBuildPipeline}
                                options={getBuildPipelineOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Build Status'
                                value={selectedBuildStatus}
                                onChange={handleSelectBuildStatus}
                                options={buildStatusOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType.includes('release') && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Release Pipeline Name'
                                value={selectedReleasePipeline}
                                onChange={handleSelectReleasePipeline}
                                options={getReleasePipelineOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType.includes('release.deployment') && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Stage Name'
                                value={selectedStageName}
                                onChange={handleSelectStageName}
                                options={getStageNameOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={selectedReleasePipeline === filterLabelValuePairAll.value || isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                (eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalPending ||
                    eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalComplete) && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Approval Type'
                                value={selectedApprovalType}
                                onChange={handleSelectApprovalType}
                                options={releaseApprovalTypeOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalComplete && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Approval Status'
                                value={selectedApprovalStatus}
                                onChange={handleSelectApprovalStatus}
                                options={releaseApprovalStatusOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentCompleted && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Status'
                                value={selectedReleaseStatus}
                                onChange={handleSelectReleaseStatus}
                                options={releaseStatusOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
        </>
    );
};

export default PipelinesFilter;
