import React from 'react';

type SvgIconNames = 'check';

const SVGIcons: Record<SvgIconNames, JSX.Element> = {
    check: (
        <>
            {/* <svg xmlns='http://www.w3.org/2000/svg' width='58' height='58' viewBox='0 0 58 58' fill='none'> */}
            <path
                d='M39.5 22.5L26.9 35.5L18.5 26.8333'
                stroke='#4273EB'
                strokeWidth='3'
            />
            <rect
                x='2'
                y='2'
                width='54'
                height='54'
                rx='27'
                stroke='#4273EB'
                strokeWidth='3'
            />
        </>
    ),
};

export default SVGIcons;
