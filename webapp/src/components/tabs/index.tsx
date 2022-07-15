import React from 'react';
import {Tabs as ReactBootstrapTabs, Tab} from 'react-bootstrap';

import './styles.scss';

type TabsProps = {
    tabs: TabData[]
    tabsClassName?: string;
}

const Tabs = ({tabs, tabsClassName = ''}: TabsProps) => {
    return (
        <ReactBootstrapTabs
            id='tab-component'
            className={`tabs ${tabsClassName}`}
        >
            {
                tabs.map((tabData, index) => (
                    <Tab
                        key={tabData.title}
                        eventKey={index}
                        title={tabData.title}
                    >
                        {tabData.tabPanel}
                    </Tab>
                ))
            }
        </ReactBootstrapTabs>
    );
};

export default Tabs;
