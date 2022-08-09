import {createSlice} from '@reduxjs/toolkit';

export interface CreateTaskModal {
    visibility: boolean
}

const initialState: CreateTaskModal = {
    visibility: false,
};

export const openTaskModalSlice = createSlice({
    name: 'openTaskModal',
    initialState,
    reducers: {
        showModal: (state) => {
            state.visibility = true;
        },
        hideModal: (state) => {
            state.visibility = false;
        },
    },
});

export const {showModal, hideModal} = openTaskModalSlice.actions;

export default openTaskModalSlice.reducer;
