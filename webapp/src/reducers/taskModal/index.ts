import {createSlice} from '@reduxjs/toolkit';

export interface CreateTaskModal {
    visibility: boolean
    title: string
    description: string
}

const initialState: CreateTaskModal = {
    visibility: false,
    title: '',
    description: '',
};

export const openTaskModalSlice = createSlice({
    name: 'openTaskModal',
    initialState,
    reducers: {
        showModal: (state, action) => {
            state.visibility = true;
            state.title = action.payload[3];
            state.description = action.payload[4];
        },
        hideModal: (state) => {
            state.visibility = false;
            state.title = '';
            state.description = '';
        },
    },
});

export const {showModal, hideModal} = openTaskModalSlice.actions;

export default openTaskModalSlice.reducer;
