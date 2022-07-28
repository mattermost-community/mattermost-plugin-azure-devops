import React from 'react';

import Tabs from 'components/tabs';

import Connections from './connections';
import Subscriptions from './subscriptions';

import './styles.scss';

const Rhs = (): JSX.Element => {
    const tabData: TabData[] = [
        {
            title: 'Boards',
            tabPanel: <div>{'Hello!'}</div>,
        },
        {
            title: 'Subscription',
            tabPanel: <Subscriptions/>,
        },
        {
            title: 'Connections',
            tabPanel: <Connections/>,
        },
    ];

    return (
        <Tabs
            tabs={tabData}
            tabsClassName='rhs-tabs'
        />
    );
};

export default Rhs;
