import {createSlice, PayloadAction} from '@reduxjs/toolkit';

type UserConnectedState = {
    isConnected: boolean;
};

const initialState: UserConnectedState = {
    isConnected: false,
};

export const userConnectedSlice = createSlice({
    name: 'userConnectedSlice',
    initialState,
    reducers: {
        toggleIsConnected: (state: UserConnectedState, action: PayloadAction<boolean>) => {
            state.isConnected = action.payload;
        },
    },
});

export const {toggleIsConnected} = userConnectedSlice.actions;

export default userConnectedSlice.reducer;
