import React from 'react';
import Tab from 'react-bootstrap/Tab';
import {Tabs as BootstrapTabs} from 'react-bootstrap';

import './styles.scss';

type TabsProps = {
    tabs: TabsData[]
}

const Tabs = ({tabs}: TabsProps) => {
    return (
        <BootstrapTabs
            defaultActiveKey='profile'
            id='RhsTabs'
            className='mb-3'
        >
            {tabs.map((tab, index) => (
                <Tab
                    key={index}
                    eventKey={tab.title}
                    title={tab.title}
                >
                    {tab.component}
                </Tab>
            ))}
        </BootstrapTabs>
    );
};

export default Tabs;
