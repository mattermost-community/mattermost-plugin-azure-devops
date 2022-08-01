import {createSlice} from '@reduxjs/toolkit';

const initialState: TaskModalState = {
    visibility: false,
};

export const openTaskModalSlice = createSlice({
    name: 'openTaskModal',
    initialState,
    reducers: {
        showTaskModal: (state) => {
            state.visibility = true;
        },
        hideTaskModal: (state) => {
            state.visibility = false;
        },
    },
});

export const {showTaskModal, hideTaskModal} = openTaskModalSlice.actions;

export default openTaskModalSlice.reducer;
