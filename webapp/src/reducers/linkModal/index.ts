import {createSlice, PayloadAction} from '@reduxjs/toolkit';

import {getProjectLinkModalArgs} from 'utils';

const initialState: LinkProjectModalState = {
    visibility: false,
    organization: '',
    project: '',
    isLinked: false,
};

export const linkProjectModalSlice = createSlice({
    name: 'linkProjectModalSlice',
    initialState,
    reducers: {
        toggleShowLinkModal: (state: LinkProjectModalState, action: PayloadAction<GlobalModalActionPayload>) => {
            state.visibility = action.payload.isVisible;
            state.organization = '';
            state.project = '';
            state.isLinked = action.payload.isActionDone ?? false;

            if (action.payload.commandArgs.length > 0) {
                const {organization, project} = getProjectLinkModalArgs(action.payload.commandArgs[0]) as LinkPayload;
                state.organization = organization;
                state.project = project;
            }
        },
        toggleIsLinkedProjectListChanged: (state: LinkProjectModalState, action: PayloadAction<boolean>) => {
            state.isLinked = action.payload;
        },
    },
});

export const {toggleShowLinkModal, toggleIsLinkedProjectListChanged} = linkProjectModalSlice.actions;

export default linkProjectModalSlice.reducer;
