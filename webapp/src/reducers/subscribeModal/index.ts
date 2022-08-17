import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
    isCreated: false,
};

export const openSubscribeModalSlice = createSlice({
    name: 'openSubscribeModal',
    initialState,
    reducers: {
        toggleShowSubscribeModal: (state: SubscribeModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.isCreated = false;
        },
        toggleIsSubscribed: (state: SubscribeModalState, action: PayloadAction<boolean>) => {
            state.isCreated = action.payload;
        },
    },
});

export const {toggleShowSubscribeModal, toggleIsSubscribed} = openSubscribeModalSlice.actions;

export default openSubscribeModalSlice.reducer;
