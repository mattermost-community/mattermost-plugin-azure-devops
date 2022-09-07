import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';

import './styles.scss';
import SVGWrapper from 'components/svgWrapper';
import plugin_constants from 'plugin_constants';

type SubscriptionCardProps = {
    handleDeleteSubscrption: (subscriptionDetails: SubscriptionDetails) => void
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {channelType, eventType, channelName}, subscriptionDetails}: SubscriptionCardProps) => (
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
                            {plugin_constants.SVGIcons.workEvent}
                        </SVGWrapper>
                    }
                    labelExtraClassName='margin-left-5'
                    value={plugin_constants.common.boardsEventTypeMap[eventType as EventType]}
                />
                <LabelValuePair
                    labelIconClassName={`icon ${channelType === plugin_constants.common.channelType.priivate ? 'icon-lock-outline' : 'icon-globe'} icon-label`}
                    value={channelName}
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
