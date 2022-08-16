import {setGlobalModalState} from 'reducers/globalModal';
import {getCommandArgs} from 'utils';

export default class Hooks {
    constructor(store) {
        this.store = store;
    }

    slashCommandWillBePostedHook = (message, contextArgs) => {
        let commandTrimmed;
        if (message) {
            commandTrimmed = message.trim();
        }

        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops link')) {
            const commandArgs = getCommandArgs(commandTrimmed);
            this.store.dispatch(setGlobalModalState({modalId: 'linkProject', commandArgs}));
            return Promise.resolve({
                message,
                args: contextArgs,
            });
        }

        if (commandTrimmed && commandTrimmed.startsWith('/azuredevops boards create')) {
            const commandArgs = getCommandArgs(commandTrimmed);
            this.store.dispatch(setGlobalModalState({modalId: 'createBoardTask', commandArgs}));
            return Promise.resolve({
                message,
                args: contextArgs,
            });
        }

        return Promise.resolve({
            message,
            args: contextArgs,
        });
    }
}
