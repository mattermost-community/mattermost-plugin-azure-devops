import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getProjectLinkDetails} from 'utils';

const initialState: LinkProjectModalState = {
    visibility: false,
    organization: '',
    project: '',
    isLinked: false,
};

export const openLinkModalSlice = createSlice({
    name: 'openLinkModal',
    initialState,
    reducers: {
        showLinkModal: (state: LinkProjectModalState, action: PayloadAction<Array<string>>) => {
            if (action.payload.length > 2) {
                const details = getProjectLinkDetails(action.payload[2]);
                if (details.length === 2) {
                    state.organization = details[0];
                    state.project = details[1];
                }
            }
            state.visibility = true;
            state.isLinked = false;
        },
        hideLinkModal: (state: LinkProjectModalState) => {
            state.visibility = false;
            state.organization = '';
            state.project = '';
        },
        toggleIsLinked: (state: LinkProjectModalState, action: PayloadAction<boolean>) => {
            state.isLinked = action.payload;
        },
    },
});

export const {showLinkModal, hideLinkModal, toggleIsLinked} = openLinkModalSlice.actions;

export default openLinkModalSlice.reducer;
