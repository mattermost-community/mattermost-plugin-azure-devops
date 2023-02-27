import React from 'react';

import mm_constants from 'mattermost-redux/constants/general';

import pluginConstants from 'pluginConstants';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';
import SVGWrapper from 'components/svgWrapper';
import Chip from 'components/chip';
import Tooltip from 'components/tooltip';

import './styles.scss';

type SubscriptionCardProps = {
    handleDeleteSubscrption: (subscriptionDetails: SubscriptionDetails) => void
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {channelType, eventType, serviceType, channelName, createdBy, targetBranch, repositoryName, pullRequestCreatedByName, pullRequestReviewersContainsName, pushedByName, mergeResultName, notificationTypeName, areaPath, releasePipelineName, buildPipeline, buildStatusName, approvalStatusName, approvalTypeName, releaseStatusName, stageNameValue, runPipelineName, runEnvironment, runStage, runStageId, runResultId, runStageResultId, runStageStateIdName, runStateIdName}, subscriptionDetails}: SubscriptionCardProps) => {
    const showFilter = areaPath || repositoryName || targetBranch || pullRequestCreatedByName || pullRequestReviewersContainsName || pushedByName || mergeResultName || notificationTypeName || releasePipelineName || buildPipeline || buildStatusName || approvalStatusName || approvalTypeName || releaseStatusName || stageNameValue || runPipelineName || runEnvironment || runStage || runStageId || runResultId || runStageResultId || runStageStateIdName || runStateIdName;

    return (
        <BaseCard>
            <>
                <div className='d-flex justify-content-between align-items-center mb-2'>
                    <div className='d-flex align-item-center'>
                        <SVGWrapper
                            width={20}
                            height={20}
                            viewBox={pluginConstants.common.serviceTypeIcon[serviceType].viewBox}
                        >
                            {pluginConstants.common.serviceTypeIcon[serviceType].icon}
                        </SVGWrapper>
                        <p className='ml-1 mb-0 font-bold text-capitalize'>{serviceType}</p>
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
                        icon={{
                            className: 'icon azure-devops-icon azure-devops-icon-event icon-event-type',
                            tooltipText: 'Event Type',
                        }}
                        value={pluginConstants.common.eventTypeMap[eventType as EventType] ?? ''}
                    />
                    <LabelValuePair
                        icon={{
                            className: `icon ${channelType === mm_constants.PRIVATE_CHANNEL ? 'icon-lock-outline' : 'icon-globe'} icon-label`,
                            tooltipText: `${channelType === mm_constants.PRIVATE_CHANNEL ? 'Private Channel' : 'Public Channel'}`,
                        }}
                        value={channelName}
                    />
                    <LabelValuePair
                        icon={{
                            className: 'icon icon-account-outline icon-label',
                            tooltipText: 'Created By',
                        }}
                        value={`Subscription created by ${createdBy}`}
                    />
                    {
                        showFilter && (
                            <div className='d-flex align-item-center margin-left-5'>
                                <div className='card-filter'>
                                    <Tooltip tooltipContent={'Filter(s)'}>
                                        <i
                                            className='azure-devops-icon azure-devops-icon-filter'
                                            aria-hidden='true'
                                        />
                                    </Tooltip>
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
                                    {(runStage || runStageId) && <Chip text={`Stage is: ${runStage || runStageId}`}/>}
                                    {runStageStateIdName && <Chip text={`Stage state is: ${runStageStateIdName}`}/>}
                                    {runStageResultId && <Chip text={`Stage result is: ${runStageResultId}`}/>}
                                    {runStateIdName && <Chip text={`State is: ${runStateIdName}`}/>}
                                    {runResultId && <Chip text={`Result is: ${runResultId}`}/>}
                                </div>
                            </div>
                        )
                    }
                </div>
            </>
        </BaseCard>
    );
};

export default SubscriptionCard;
