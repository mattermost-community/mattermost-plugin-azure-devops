import React from 'react';

import BaseCard from 'components/card/base';
import IconButton from 'components/buttons/iconButton';
import LabelValuePair from 'components/labelValuePair';

import './styles.scss';

type SubscriptionCardProps = {
    subscriptionDetails: SubscriptionDetails
}

const SubscriptionCard = ({subscriptionDetails: {id, name, eventType}}: SubscriptionCardProps) => {
    return (
        <BaseCard>
            <div className='d-flex'>
                <div className='project-details'>
                    <p className='margin-bottom-10 font-size-14 font-bold'>{id}</p>
                    <LabelValuePair
                        label='Name'
                        value={name}
                    />
                    <LabelValuePair
                        label='Event'
                        value={eventType}
                    />
                </div>
                <div className='button-wrapper'>
                    <IconButton
                        tooltipText='Delete subscription'
                        iconClassName='fa fa-trash-o'
                        extraClass='delete-button'
                    />
                </div>
            </div>
        </BaseCard>
    );
};

export default SubscriptionCard;
