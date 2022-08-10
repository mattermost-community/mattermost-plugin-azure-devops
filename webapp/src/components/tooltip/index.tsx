import React from 'react';
import {Placement} from 'react-bootstrap/esm/types';

import {OverlayTrigger, Tooltip as ReactTooltip} from 'react-bootstrap';

// Styles
import './styles.scss';

type TooltipProps = {
    tooltipContent: JSX.Element | string
    children: JSX.Element
    placement?: Placement
}

const Tooltip = ({tooltipContent, children, placement = 'top'}: TooltipProps) => {
    return (
        <OverlayTrigger
            placement={placement}
            overlay={
                <ReactTooltip
                    id='tooltip'
                    className='tooltip-wrapper'
                    placement={placement}
                >
                    {tooltipContent}
                </ReactTooltip>
            }
        >
            {children}
        </OverlayTrigger>
    );
};

export default Tooltip;
