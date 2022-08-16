import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: SubscribeModalState = {
    visibility: false,
};

export const openSubscribeModalSlice = createSlice({
    name: 'openSubscribeModal',
    initialState,
    reducers: {
        toggleShowSubscribeModal: (state: SubscribeModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
        },
    },
});

export const {toggleShowSubscribeModal} = openSubscribeModalSlice.actions;

export default openSubscribeModalSlice.reducer;
