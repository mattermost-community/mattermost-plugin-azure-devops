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

            if (action.payload.args?.length === 2) {
                state.organization = action.payload.args[0];
                state.project = action.payload.args[1];
                return;
            }

            state.organization = null;
            state.project = null;
        },
        toggleIsSubscribed: (state: SubscribeModalState, action: PayloadAction<boolean>) => {
            state.isCreated = action.payload;
        },
    },
});

export const {toggleShowSubscribeModal, toggleIsSubscribed} = subscriptionModalSlice.actions;

export default subscriptionModalSlice.reducer;
