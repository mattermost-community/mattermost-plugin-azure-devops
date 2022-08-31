import React from 'react';

import {useDispatch} from 'react-redux';

import EmptyState from 'components/emptyState';

import Utils from 'utils';

const AccountNotLinked = () => {
    const dispatch = useDispatch();

    const closeRHS = () => {
        dispatch({
            type: 'UPDATE_RHS_STATE',
            state: null,
        });
    };

    // Opens link project modal
    const handleConnectAccount = () => {
        window.open(`${Utils.getBaseUrls().pluginApiBaseUrl}/oauth/connect`, '_blank');
        closeRHS();
    };

    return (
        <>
            <EmptyState
                title='No account connected'
                subTitle={{text: 'Connect your account by clicking the button below'}}
                buttonText='Connect your account'
                buttonAction={handleConnectAccount}
                icon='azure'
            />
        </>
    );
};

export default AccountNotLinked;
