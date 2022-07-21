import {showModal} from 'reducers/taskModal';
import {splitArgs} from '../utils';

export default class Hooks {
    constructor(store) {
        this.store = store;
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

        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops boards create')) {
            const args = splitArgs(commandTrimmed);
            this.store.dispatch(showModal(args));
            return Promise.resolve({});
        }
        return Promise.resolve({
            message,
            args: contextArgs,
        });
    }
}
