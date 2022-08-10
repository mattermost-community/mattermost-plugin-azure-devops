import {createSlice} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
};

export const openSubscribeModalSlice = createSlice({
    name: 'openSubscribeModal',
    initialState,
    reducers: {
        showSubscribeModal: (state: SubscribeModalState) => {
            state.visibility = true;
        },
        hideSubscribeModal: (state: SubscribeModalState) => {
            state.visibility = false;
        },
    },
});

export const {showSubscribeModal, hideSubscribeModal} = openSubscribeModalSlice.actions;

export default openSubscribeModalSlice.reducer;
