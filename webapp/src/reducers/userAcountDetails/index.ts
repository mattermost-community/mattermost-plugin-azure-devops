import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getProjectLinkDetails} from 'utils';

const initialState: userConnectionState = {
    isConnectionTriggered: false,
    isUserDisconnected: true,
};

export const UserConnectionSlice = createSlice({
    name: 'userAccountDetails',
    initialState,
    reducers: {
        toggleConnectionTriggered: (state: userConnectionState, action: PayloadAction<boolean>) => {
            state.isConnectionTriggered = action.payload;
        },
        toggleIsDisconnected: (state: userConnectionState, action: PayloadAction<boolean>) => {
            state.isUserDisconnected = action.payload;
        },
    },
});

export const {toggleConnectionTriggered, toggleIsDisconnected} = UserConnectionSlice.actions;

export default UserConnectionSlice.reducer;
