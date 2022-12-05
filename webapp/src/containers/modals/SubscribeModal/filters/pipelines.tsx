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
    selectedRunPipeline: string
    handleSelectRunPipeline: (value: string, name?: string) => void
    selectedRunStage: string
    handleSelectRunStage: (value: string, name?: string) => void
    selectedRunEnvironment: string
    handleSelectRunEnvironment: (value: string, name?: string) => void
    selectedRunStageId: string
    handleSelectRunStageId: (value: string, name?: string) => void
    selectedRunStageStateId: string
    handleSelectRunStageStateId: (value: string, name?: string) => void
    selectedRunStageResultId: string
    handleSelectRunStageResultId: (value: string, name?: string) => void
    selectedRunStateId: string
    handleSelectRunStateId: (value: string, name?: string) => void
    selectedRunResultId: string
    handleSelectRunResultId: (value: string, name?: string) => void
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
    handleSelectRunPipeline,
    selectedRunPipeline,
    handleSelectRunStage,
    selectedRunStage,
    handleSelectRunEnvironment,
    selectedRunEnvironment,
    selectedRunStageId,
    handleSelectRunStageId,
    selectedRunStageStateId,
    handleSelectRunStageStateId,
    selectedRunStageResultId,
    handleSelectRunStageResultId,
    selectedRunStateId,
    handleSelectRunStateId,
    selectedRunResultId,
    handleSelectRunResultId,
}: PipelinesFilterProps) => {
    const {buildStatusOptions, releaseApprovalTypeOptions, releaseApprovalStatusOptions, releaseStatusOptions, subscriptionFiltersForPipelines, subscriptionFiltersNameForPipelines, runStageStateIdOptions, runStageResultIdOptions, runStateIdOptions, runResultIdOptions} = pluginConstants.form;

    const {getApiState, makeApiRequestWithCompletionStatus} = usePluginApi();

    const getSubscriptionFiltersRequest = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForPipelines,
        eventType,
        releasePipelineId: selectedReleasePipeline,
        runPipeline: selectedRunPipeline,
    }), [organization, projectId, eventType, subscriptionFiltersForPipelines, selectedBuildPipeline, selectedReleasePipeline, selectedRunPipeline]);

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
            return;
        }

        setIsFiltersError(false);
    }, [isLoading, isError, isSuccess]);

    const getBuildPipelineOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.buildPipeline], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getReleasePipelineOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.releasePipelineName], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getStageNameOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.stageName], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getRunPipelineOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.runPipeline], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getRunStageOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.runStage], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getRunEnvironmentOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.runEnvironment], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);
    const getRunStageIdOptions = useCallback(() => (isSuccess ? ([{...filterLabelValuePairAll}, ...formLabelValuePairs('displayValue', 'value', filtersData[subscriptionFiltersNameForPipelines.runStageId], ['[Any]'])]) : [pluginConstants.common.filterLabelValuePairAll]), [filtersData]);

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
                )
            }
            {
                eventType.includes('release.deployment') && (
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
                )
            }
            {
                (eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalPending ||
                    eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalComplete) && (
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
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentApprovalComplete && (
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
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.releaseDeploymentCompleted && (
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
                )
            }
            {
                (eventType === pluginConstants.common.eventTypePipelineKeys.runStageApprovalComplete || eventType === pluginConstants.common.eventTypePipelineKeys.runStageApprovalPending || eventType === pluginConstants.common.eventTypePipelineKeys.runStageStateChanged || eventType === pluginConstants.common.eventTypePipelineKeys.runStateChanged) && (
                    <div className='margin-bottom-10'>
                        <Dropdown
                            placeholder='Pipeline'
                            value={selectedRunPipeline}
                            onChange={handleSelectRunPipeline}
                            options={getRunPipelineOptions()}
                            error={isError}
                            loadingOptions={isLoading}
                            disabled={isLoading}
                        />
                    </div>
                )
            }
            {
                (eventType === pluginConstants.common.eventTypePipelineKeys.runStageApprovalComplete || eventType === pluginConstants.common.eventTypePipelineKeys.runStageApprovalPending) && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Stage'
                                value={selectedRunStage}
                                onChange={handleSelectRunStage}
                                options={getRunStageOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={selectedRunPipeline === filterLabelValuePairAll.value || isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Environment'
                                value={selectedRunEnvironment}
                                onChange={handleSelectRunEnvironment}
                                options={getRunEnvironmentOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.runStageStateChanged && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Stage'
                                value={selectedRunStageId}
                                onChange={handleSelectRunStageId}
                                options={getRunStageIdOptions()}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={selectedRunPipeline === filterLabelValuePairAll.value || isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='State'
                                value={selectedRunStageStateId}
                                onChange={handleSelectRunStageStateId}
                                options={runStageStateIdOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Result'
                                value={selectedRunStageResultId}
                                onChange={handleSelectRunStageResultId}
                                options={runStageResultIdOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={(selectedRunStageStateId !== filterLabelValuePairAll.value && selectedRunStageStateId !== 'Completed') || isLoading}
                            />
                        </div>
                    </>
                )
            }
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.runStateChanged && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='State'
                                value={selectedRunStateId}
                                onChange={handleSelectRunStateId}
                                options={runStateIdOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Result'
                                value={selectedRunResultId}
                                onChange={handleSelectRunResultId}
                                options={runResultIdOptions}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={(selectedRunStateId !== filterLabelValuePairAll.value && selectedRunStateId !== 'Completed') || isLoading}
                            />
                        </div>
                    </>
                )
            }
        </>
    );
};

export default PipelinesFilter;
