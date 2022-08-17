import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getCreateTaskModalCommandArgs} from 'utils';

const initialState: CreateTaskModalState = {
    visibility: false,
    commandArgs: {
        title: '',
        description: '',
    },
};

export const openTaskModalSlice = createSlice({
    name: 'openTaskModal',
    initialState,
    reducers: {
        toggleShowTaskModal: (state: CreateTaskModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.commandArgs.title = '';
            state.commandArgs.description = '';

            if (action.payload.commandArgs.length > 1) {
                const {title, description} = getCreateTaskModalCommandArgs(action.payload.commandArgs) as TaskFieldsCommandArgs;
                state.commandArgs.title = title;
                state.commandArgs.description = description;
            }
        },
    },
});

export const {toggleShowTaskModal} = openTaskModalSlice.actions;

export default openTaskModalSlice.reducer;
