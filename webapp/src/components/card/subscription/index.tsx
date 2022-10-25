import React from 'react';

import mm_constants from 'mattermost-redux/constants/general';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';
import SVGWrapper from 'components/svgWrapper';

import pluginConstants from 'pluginConstants';

import './styles.scss';

type SubscriptionCardProps = {
    handleDeleteSubscrption: (subscriptionDetails: SubscriptionDetails) => void
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {channelType, eventType, serviceType, channelName, createdBy}, subscriptionDetails}: SubscriptionCardProps) => (
    <BaseCard>
        <>
            <div className='d-flex justify-content-between align-items-center mb-2'>
                <div className='d-flex'>
                    <SVGWrapper
                        width={20}
                        height={20}
                        viewBox=' 0 0 20 20'
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
            </div>
        </>
    </BaseCard>
);

export default SubscriptionCard;
