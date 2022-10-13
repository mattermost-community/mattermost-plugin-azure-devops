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
        <div className='d-flex'>
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
                {/* TODO: Modify UI to display service type according to the new designs */}
                <LabelValuePair
                    value={pluginConstants.common.serviceTypeMap[eventType as EventType]}
                />
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
    </BaseCard>
);

export default SubscriptionCard;
