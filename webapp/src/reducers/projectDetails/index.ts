import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: ProjectDetails = {
    mattermostID: '',
    projectID: '',
    projectName: '',
    organizationName: '',
};

export const projectDetailsSlice = createSlice({
    name: 'projectDetails',
    initialState,
    reducers: {
        setProjectDetails: (state: ProjectDetails, action: PayloadAction<ProjectDetails>) => {
            state.mattermostID = action.payload.mattermostID;
            state.projectID = action.payload.projectID;
            state.projectName = action.payload.projectName;
            state.organizationName = action.payload.organizationName;
        },
        resetProjectDetails: (state: ProjectDetails) => {
            state.mattermostID = '';
            state.projectID = '';
            state.projectName = '';
            state.organizationName = '';
        },
    },
});

// Action creators are generated for each case reducer function
export const {setProjectDetails, resetProjectDetails} = projectDetailsSlice.actions;

export default projectDetailsSlice.reducer;
