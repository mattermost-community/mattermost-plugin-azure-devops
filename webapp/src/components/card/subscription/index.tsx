import React from 'react';

import mm_constants from 'mattermost-redux/constants/general';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';
import SVGWrapper from 'components/svgWrapper';

import pluginConstants from 'pluginConstants';

import './styles.scss';
import Chip from 'components/chip';

type SubscriptionCardProps = {
    handleDeleteSubscrption: (subscriptionDetails: SubscriptionDetails) => void
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {channelType, eventType, serviceType, channelName, createdBy, targetBranch, repositoryName, pullRequestCreatedByName, pullRequestReviewersContainsName, pushedByName, mergeResultName, notificationTypeName, areaPath, releasePipelineName, buildPipeline, buildStatusName, approvalStatusName, approvalTypeName, releaseStatusName, stageNameValue, runPipelineName, runEnvironment, runStage, runStageId, runResultId, runStageResultId, runStageStateIdName, runStateIdName}, subscriptionDetails}: SubscriptionCardProps) => {
    let serviceTypeIcon = pluginConstants.SVGIcons.boards;
    if (serviceType === pluginConstants.common.repos) {
        serviceTypeIcon = pluginConstants.SVGIcons.repos;
    } else if (serviceType === pluginConstants.common.pipelines) {
        serviceTypeIcon = pluginConstants.SVGIcons.pipelines;
    }

    return (
        <BaseCard>
            <div>
                <div className='d-flex justify-content-between align-items-center mb-2'>
                    <div className='d-flex align-item-center'>
                        <SVGWrapper
                            width={20}
                            height={20}
                            viewBox={serviceType === pluginConstants.common.pipelines ? ' 0 0 17 17' : ' 0 0 16 16'}
                        >
                            {serviceTypeIcon}
                        </SVGWrapper>
                        <p className={`ml-1 mb-0 font-bold color-${serviceType} text-capitalize`}>{serviceType}</p>
                    </div>
                    <div className='button-wrapper'>
                        <IconButton
                            tooltipText='Delete subscription'
                            iconClassName='fa fa-trash-o'
                            extraClass='delete-button'
                            onClick={() => handleDeleteSubscrption(subscriptionDetails)}
                        />
                    </div>
                </div>
                <div className='project-details'>
                    <LabelValuePair
                        label={
                            <SVGWrapper
                                width={12}
                                height={12}
                                viewBox='0 0 10 10'
                            >
                                {pluginConstants.SVGIcons.workEvent}
                            </SVGWrapper>
                        }
                        labelExtraClassName='margin-left-5'
                        value={pluginConstants.common.eventTypeMap[eventType as EventType] ?? ''}
                    />
                    <LabelValuePair
                        labelIconClassName={`icon ${channelType === mm_constants.PRIVATE_CHANNEL ? 'icon-lock-outline' : 'icon-globe'} icon-label`}
                        value={channelName}
                    />
                    <LabelValuePair
                        labelIconClassName={'icon icon-account-outline icon-label'}
                        value={`Subscription created by ${createdBy}`}
                    />
                    {
                        (areaPath || repositoryName || targetBranch || pullRequestCreatedByName || pullRequestReviewersContainsName || pushedByName || mergeResultName || notificationTypeName || releasePipelineName || buildPipeline || buildStatusName || approvalStatusName || approvalTypeName || releaseStatusName || stageNameValue || runPipelineName || runEnvironment || runStage || runStageId || runResultId || runStageResultId || runStageStateIdName || runStateIdName) && (
                            <div className='d-flex align-item-center margin-left-5'>
                                <div className='card-filter'>
                                    <SVGWrapper
                                        width={14}
                                        height={14}
                                        viewBox='0 0 12 12'
                                    >
                                        {pluginConstants.SVGIcons.cardFilter}
                                    </SVGWrapper>
                                </div>
                                <div className='card-chip-wrapper'>
                                    {

                                        // Remove the extra character "/" from start and end of the area path string returned by the API
                                        areaPath && <Chip text={`Area path - ${areaPath.substring(1, areaPath.length - 1)}`}/>
                                    }
                                    {repositoryName && <Chip text={`Repository is: ${repositoryName}`}/>}
                                    {targetBranch && <Chip text={`Target branch is: ${targetBranch}`}/>}
                                    {pullRequestCreatedByName && <Chip text={`Requested by a member of group: ${pullRequestCreatedByName}`}/>}
                                    {pullRequestReviewersContainsName && <Chip text={`Reviewer includes group: ${pullRequestReviewersContainsName}`}/>}
                                    {pushedByName && <Chip text={`Pushed by a member of group: ${pushedByName}`}/>}
                                    {mergeResultName && <Chip text={`Merge result: ${mergeResultName}`}/>}
                                    {notificationTypeName && <Chip text={`Change: ${notificationTypeName}`}/>}
                                    {releasePipelineName && <Chip text={`Release pipeline is: ${releasePipelineName}`}/>}
                                    {buildPipeline && <Chip text={`Build pipeline is: ${buildPipeline}`}/>}
                                    {buildStatusName && <Chip text={`Build status is: ${buildStatusName}`}/>}
                                    {stageNameValue && <Chip text={`Stage name is: ${stageNameValue}`}/>}
                                    {approvalStatusName && <Chip text={`Approval status is: ${approvalStatusName}`}/>}
                                    {approvalTypeName && <Chip text={`Approval type is: ${approvalTypeName}`}/>}
                                    {releaseStatusName && <Chip text={`Release status is: ${releaseStatusName}`}/>}
                                    {runPipelineName && <Chip text={`Pipeline is: ${runPipelineName}`}/>}
                                    {runEnvironment && <Chip text={`Environment is: ${runEnvironment}`}/>}
                                    {runStage && <Chip text={`Stage is: ${runStage}`}/>}
                                    {runStageId && <Chip text={`Stage is: ${runStageId}`}/>}
                                    {runStageStateIdName && <Chip text={`Stage state is: ${runStageStateIdName}`}/>}
                                    {runStageResultId && <Chip text={`Stage result is: ${runStageResultId}`}/>}
                                    {runStateIdName && <Chip text={`State is: ${runStateIdName}`}/>}
                                    {runResultId && <Chip text={`Result is: ${runResultId}`}/>}
                                </div>
                            </div>
                        )
                    }
                </div>
            </div>
        </BaseCard>
    );
};

export default SubscriptionCard;
