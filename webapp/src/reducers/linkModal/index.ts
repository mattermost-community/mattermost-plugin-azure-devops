import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getProjectLinkModalArgs} from 'utils';

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
        toggleShowLinkModal: (state: LinkProjectModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.organization = '';
            state.project = '';
            state.isLinked = false;

            if (action.payload.commandArgs.length > 0) {
                const {organization, project} = getProjectLinkModalArgs(action.payload.commandArgs[0]) as LinkPayload;
                state.organization = organization;
                state.project = project;
            }
        },
        toggleIsLinked: (state: LinkProjectModalState, action: PayloadAction<boolean>) => {
            state.isLinked = action.payload;
        },
    },
});

export const {toggleShowLinkModal, toggleIsLinked} = openLinkModalSlice.actions;

export default openLinkModalSlice.reducer;
