import React, {useMemo} from 'react';

import pluginConstants from 'pluginConstants';
import {filterLabelValuePairAll} from 'pluginConstants/common';

import Dropdown from 'components/dropdown';

import useLoadFilters from 'hooks/useLoadFilters';

type PipelinesFilterProps = {
    organization: string
    projectId: string
    eventType: string
    selectedBuildPipeline: string
    isModalOpen: boolean
    handleSetFilter: HandleSetSubscriptionFilter
    setIsFiltersError: (value: boolean) => void
    selectedBuildStatus: string
    selectedReleasePipeline: string
    selectedStageName: string
    selectedApprovalType: string
    selectedApprovalStatus: string
    selectedReleaseStatus: string
    selectedRunPipeline: string
    selectedRunStage: string
    selectedRunEnvironment: string
    selectedRunStageId: string
    selectedRunStageStateId: string
    selectedRunStageResultId: string
    selectedRunStateId: string
    selectedRunResultId: string
}

const PipelinesFilter = ({
    organization,
    projectId,
    eventType,
    selectedBuildPipeline,
    isModalOpen,
    handleSetFilter,
    setIsFiltersError,
    selectedBuildStatus,
    selectedReleasePipeline,
    selectedStageName,
    selectedApprovalType,
    selectedApprovalStatus,
    selectedReleaseStatus,
    selectedRunPipeline,
    selectedRunStage,
    selectedRunEnvironment,
    selectedRunStageId,
    selectedRunStageStateId,
    selectedRunStageResultId,
    selectedRunStateId,
    selectedRunResultId,
}: PipelinesFilterProps) => {
    const {buildStatusOptions, releaseApprovalTypeOptions, releaseApprovalStatusOptions, releaseStatusOptions, subscriptionFiltersForPipelines, subscriptionFiltersNameForPipelines, runStageStateIdOptions, runStageResultIdOptions, runStateIdOptions, runResultIdOptions} = pluginConstants.form;

    const getSubscriptionFiltersRequestParams = useMemo<GetSubscriptionFiltersRequest>(() => ({
        organization,
        projectId,
        filters: subscriptionFiltersForPipelines,
        eventType,
        releasePipelineId: selectedReleasePipeline,
        runPipeline: selectedRunPipeline,
    }), [organization, projectId, eventType, subscriptionFiltersForPipelines, selectedBuildPipeline, selectedReleasePipeline, selectedRunPipeline]);

    const {filtersData, isError, isLoading, getFilterOptions} = useLoadFilters({isModalOpen, getSubscriptionFiltersRequestParams, setIsFiltersError});

    return (
        <>
            {
                eventType === pluginConstants.common.eventTypePipelineKeys.buildCompleted && (
                    <>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Build Pipeline'
                                value={selectedBuildPipeline}
                                onChange={(newValue) => handleSetFilter('buildPipeline', newValue)}
                                options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.buildPipeline])}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Build Status'
                                value={selectedBuildStatus}
                                onChange={(newValue, label) => handleSetFilter('buildStatus', newValue, 'buildStatusName', label)}
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
                            onChange={(newValue, label) => handleSetFilter('releasePipeline', newValue, 'releasePipelineName', label)}
                            options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.releasePipelineName])}
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
                            onChange={(newValue, label) => handleSetFilter('stageName', newValue, 'stageNameValue', label)}
                            options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.stageName])}
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
                            onChange={(newValue, label) => handleSetFilter('approvalType', newValue, 'approvalTypeName', label)}
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
                            onChange={(newValue, label) => handleSetFilter('approvalStatus', newValue, 'approvalStatusName', label)}
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
                            onChange={(newValue, label) => handleSetFilter('releaseStatus', newValue, 'releaseStatusName', label)}
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
                            onChange={(newValue, label) => handleSetFilter('runPipeline', newValue, 'runPipelineName', label)}
                            options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.runPipeline])}
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
                                onChange={(newValue) => handleSetFilter('runStage', newValue)}
                                options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.runStage])}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={selectedRunPipeline === filterLabelValuePairAll.value || isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='Environment'
                                value={selectedRunEnvironment}
                                onChange={(newValue) => handleSetFilter('runEnvironment', newValue)}
                                options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.runEnvironment])}
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
                                onChange={(newValue) => handleSetFilter('runStageId', newValue)}
                                options={getFilterOptions(filtersData[subscriptionFiltersNameForPipelines.runStageId])}
                                error={isError}
                                loadingOptions={isLoading}
                                disabled={selectedRunPipeline === filterLabelValuePairAll.value || isLoading}
                            />
                        </div>
                        <div className='margin-bottom-10'>
                            <Dropdown
                                placeholder='State'
                                value={selectedRunStageStateId}
                                onChange={(newValue, label) => handleSetFilter('runStageStateId', newValue, 'runStageStateIdName', label)}
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
                                onChange={(newValue) => handleSetFilter('runStageResultId', newValue)}
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
                                onChange={(newValue, label) => handleSetFilter('runStateId', newValue, 'runStateIdName', label)}
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
                                onChange={(newValue) => handleSetFilter('runResultId', newValue)}
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
