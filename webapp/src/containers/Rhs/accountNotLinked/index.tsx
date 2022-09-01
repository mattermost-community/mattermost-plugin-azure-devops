import React from 'react';

import EmptyState from 'components/emptyState';

import Utils from 'utils';

const AccountNotLinked = () => {
    // Opens link project modal
    const handleConnectAccount = () => {
        window.open(`${Utils.getBaseUrls().pluginApiBaseUrl}/oauth/connect`, '_blank');
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
