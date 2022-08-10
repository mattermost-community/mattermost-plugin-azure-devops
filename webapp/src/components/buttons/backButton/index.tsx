import React from 'react';

import {onPressingEnterKey} from 'utils';

type BackButtonProps = {
    onClick: () => void
}

const BackButton = ({onClick}: BackButtonProps) => {
    return (
        <p className='margin-bottom-25'>
            <span
                tabIndex={0}
                className='link-title'
                onKeyDown={() => onPressingEnterKey(event, onClick)}
                onClick={onClick}
            >
                <i
                    className='fa fa-caret-left'
                    aria-hidden='true'
                /> {'Back'}
            </span>
        </p>
    );
};

export default BackButton;
