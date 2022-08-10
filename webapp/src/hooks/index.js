import {showLinkModal} from 'reducers/linkModal';
import {showSubscribeModal} from 'reducers/subscribeModal';
import {splitArgs} from '../utils';

export default class Hooks {
    constructor(store) {
        this.store = store;
    }

    closeRhs() {
        this.store.dispatch({
            type: 'UPDATE_RHS_STATE',
            state: null,
        });
    }

    disconnectUser() {
        this.store.dispatch({
            type: 'userAccountDetails/toggleIsDisconnected',
            payload: true,
        });
    }

    toggleUserConnection(triggerConnection) {
        this.store.dispatch({
            type: 'userAccountDetails/toggleConnectionTriggered',
            payload: triggerConnection,
        });
    }

    slashCommandWillBePostedHook = (message, contextArgs) => {
        let commandTrimmed;
        if (message) {
            commandTrimmed = message.trim();
        }

        if (!commandTrimmed?.startsWith('/azuredevops')) {
            return Promise.resolve({
                message,
                args: contextArgs,
            });
        }

        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops link')) {
            const args = splitArgs(commandTrimmed);
            this.toggleUserConnection(true);
            this.store.dispatch(showLinkModal(args));
            return {
                message,
                args: contextArgs,
            };
        }
        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops subscribe')) {
            const args = splitArgs(commandTrimmed);
            this.toggleUserConnection(true);
            this.store.dispatch(showSubscribeModal(args));
            return {
                message,
                args: contextArgs,
            };
        }
        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops connect')) {
            this.closeRhs();
            return {
                message,
                args: contextArgs,
            };
        }
        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops disconnect')) {
            this.disconnectUser();
            this.closeRhs();
            return {
                message,
                args: contextArgs,
            };
        }
        return Promise.resolve({
            message,
            args: contextArgs,
        });
    }
}
