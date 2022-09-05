import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';

import './styles.scss';

type SubscriptionCardProps = {
    handleDeleteSubscrption: (subscriptionDetails: SubscriptionDetails) => void
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({handleDeleteSubscrption, subscriptionDetails: {projectName, eventType, channelName}, subscriptionDetails}: SubscriptionCardProps) => (
        <BaseCard>
            <div className='d-flex'>
                <div className='project-details'>
                    {/* TODO: discuss if this is required and it should be the subscription's name not the project name */}
                    {/* <LabelValuePair
                        label='Name'
                        value={projectName}
                    /> */}
                    <LabelValuePair
                        label='Event'
                        value={eventType}
                    />
                    <LabelValuePair
                        label='Channel'
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
