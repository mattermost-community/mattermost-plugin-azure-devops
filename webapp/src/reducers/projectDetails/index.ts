import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: ProjectDetails = {
    id: '',
    organization: '',
    title: '',
};

export const projectDetailsSlice = createSlice({
    name: 'projectDetails',
    initialState,
    reducers: {
        setProjectDetails: (state: ProjectDetails, action: PayloadAction<ProjectDetails>) => {
            state.id = action.payload.id;
            state.title = action.payload.title;
            state.organization = action.payload.organization;
        },
        resetProjectDetails: (state: ProjectDetails) => {
            state.id = '';
            state.title = '';
            state.organization = '';
        },
    },
});

// Action creators are generated for each case reducer function
export const {setProjectDetails, resetProjectDetails} = projectDetailsSlice.actions;

export default projectDetailsSlice.reducer;
