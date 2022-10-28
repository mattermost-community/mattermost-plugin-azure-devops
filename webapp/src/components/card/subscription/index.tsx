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

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {channelType, eventType, serviceType, channelName, createdBy, repository, targetBranch, repositoryName}, subscriptionDetails}: SubscriptionCardProps) => (
    <BaseCard>
        <>
            <div className='d-flex justify-content-between align-items-center mb-2'>
                <div className='d-flex align-item-center'>
                    <SVGWrapper
                        width={20}
                        height={20}
                        viewBox=' 0 0 16 16'
                    >
                        {serviceType === pluginConstants.common.boards ? pluginConstants.SVGIcons.boards : pluginConstants.SVGIcons.repos}
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
                    (repositoryName || targetBranch) && (
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
                                    repositoryName && <Chip text={`Repository is ${repositoryName}`}/>
                                }
                                {
                                    targetBranch && <Chip text={`Target Branch is ${targetBranch}`}/>
                                }
                            </div>
                        </div>
                    )
                }
            </div>
        </>
    </BaseCard>
);

export default SubscriptionCard;
