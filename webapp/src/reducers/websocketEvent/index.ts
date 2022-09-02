import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: WebsocketEventState = {
    isConnected: false,
    isSubscriptionDeleted: false,
};

export const websocketEventSlice = createSlice({
    name: 'websocketEventSlice',
    initialState,
    reducers: {
        toggleIsConnected: (state: WebsocketEventState, action: PayloadAction<boolean>) => {
            state.isConnected = action.payload;
        },
        toggleIsSubscriptionDeleted: (state: WebsocketEventState, action: PayloadAction<boolean>) => {
            state.isSubscriptionDeleted = action.payload;
        },
    },
});

export const {toggleIsConnected, toggleIsSubscriptionDeleted} = websocketEventSlice.actions;

export default websocketEventSlice.reducer;
