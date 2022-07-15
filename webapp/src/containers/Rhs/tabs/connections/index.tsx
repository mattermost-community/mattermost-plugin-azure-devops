/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';

import NoData from '../no_data';

const Connections = (): JSX.Element => {
    return (
        <NoData
            title='No connections yet'
            subTitle='Lorem ipsum dolor sit amet, consectetur adipiscing elit. Adipiscing nulla in tellus est mauris et eros.'
            buttonText='Create a connection'
            buttonAction={() => ''}
        />
    );
};

export default Connections;
