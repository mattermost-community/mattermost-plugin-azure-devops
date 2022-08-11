import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getProjectLinkDetails} from 'utils';

export interface CreateTaskModal {
    visibility: boolean,
    organization: string,
    project: string,
}

const initialState: CreateTaskModal = {
    visibility: false,
    organization: '',
    project: '',
};

export const openLinkModalSlice = createSlice({
    name: 'openLinkModal',
    initialState,
    reducers: {
        showLinkModal: (state: CreateTaskModal, action: PayloadAction<Array<string>>) => {
            if (action.payload.length > 2) {
                const details = getProjectLinkDetails(action.payload[2]);
                if (details.length === 2) {
                    state.organization = details[0];
                    state.project = details[1];
                }
            }
            state.visibility = true;
        },
        hideLinkModal: (state: CreateTaskModal) => {
            state.visibility = false;
            state.organization = '';
            state.project = '';
        },
    },
});

export const {showLinkModal, hideLinkModal} = openLinkModalSlice.actions;

export default openLinkModalSlice.reducer;
