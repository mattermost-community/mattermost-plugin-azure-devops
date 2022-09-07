import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
    isCreated: false,
};

export const subscriptionModalSlice = createSlice({
    name: 'subscriptionModalSlice',
    initialState,
    reducers: {
        toggleShowSubscribeModal: (state: SubscribeModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.isCreated = action.payload.isActionDone ?? false;
        },
        toggleIsSubscribed: (state: SubscribeModalState, action: PayloadAction<boolean>) => {
            state.isCreated = action.payload;
        },
    },
});

export const {toggleShowSubscribeModal, toggleIsSubscribed} = subscriptionModalSlice.actions;

export default subscriptionModalSlice.reducer;
