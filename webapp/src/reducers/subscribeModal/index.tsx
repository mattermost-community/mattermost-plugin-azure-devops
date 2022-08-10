import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
    isLinked: false,
};

export const openSubscribeModalSlice = createSlice({
    name: 'openSubscribeModal',
    initialState,
    reducers: {
        showSubscribeModal: (state: SubscribeModalState) => {
            state.visibility = true;
            state.isLinked = false;
        },
        hideSubscribeModal: (state: SubscribeModalState) => {
            state.visibility = false;
        },
        toggleIsLinked: (state: SubscribeModalState, action: PayloadAction<boolean>) => {
            state.isLinked = action.payload;
        },
    },
});

export const {showSubscribeModal, hideSubscribeModal, toggleIsLinked} = openSubscribeModalSlice.actions;

export default openSubscribeModalSlice.reducer;
