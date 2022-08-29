import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: GlobalModalState = {
    modalId: null,
    commandArgs: [],
};

export const globalModalSlice = createSlice({
    name: 'globalModalSlice',
    initialState,
    reducers: {
        setGlobalModalState: (state: GlobalModalState, action: PayloadAction<GlobalModalState>) => {
            state.modalId = action.payload.modalId;
            state.commandArgs = action.payload.commandArgs;
        },
        resetGlobalModalState: (state: GlobalModalState) => {
            state.modalId = null;
            state.commandArgs = [];
        },
    },
});

export const {setGlobalModalState, resetGlobalModalState} = globalModalSlice.actions;

export default globalModalSlice.reducer;
