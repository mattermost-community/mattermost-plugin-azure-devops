import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
    isCreated: false,
    serviceType: 'boards',
};

export const subscriptionModalSlice = createSlice({
    name: 'subscriptionModalSlice',
    initialState,
    reducers: {
        toggleShowSubscribeModal: (state: SubscribeModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.isCreated = action.payload.isActionDone ?? false;

            if (action.payload.commandArgs?.length === 3) {
                state.serviceType = action.payload.commandArgs[2];
                return;
            }

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
        setServiceType: (state: SubscribeModalState, action: PayloadAction<string>) => {
            state.serviceType = action.payload;
        },
    },
});

export const {toggleShowSubscribeModal, toggleIsSubscribed, setServiceType} = subscriptionModalSlice.actions;

export default subscriptionModalSlice.reducer;
