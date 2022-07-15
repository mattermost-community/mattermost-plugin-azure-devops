/* eslint-disable @typescript-eslint/no-unused-vars */

import React from 'react';

import NoData from '../no_data';

const Subscriptions = (): JSX.Element => {
    return (
        <NoData
            title='No subscriptions yet'
            subTitle='Lorem ipsum dolor sit amet, consectetur adipiscing elit. Adipiscing nulla in tellus est mauris et eros.'
            buttonText='Create a subscription'
            buttonAction={() => ''}
        />
    );
};

export default Subscriptions;
