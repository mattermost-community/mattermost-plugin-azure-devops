import {createSlice, PayloadAction} from '@reduxjs/toolkit';

const initialState: ProjectDetails = {
    mattermostUserID: '',
    projectID: '',
    projectName: '',
    organizationName: '',
};

export const projectDetailsSlice = createSlice({
    name: 'projectDetailsSlice',
    initialState,
    reducers: {
        setProjectDetails: (state: ProjectDetails, action: PayloadAction<ProjectDetails>) => {
            state.mattermostUserID = action.payload.mattermostUserID;
            state.projectID = action.payload.projectID;
            state.projectName = action.payload.projectName;
            state.organizationName = action.payload.organizationName;
        },
        resetProjectDetails: (state: ProjectDetails) => {
            state.mattermostUserID = '';
            state.projectID = '';
            state.projectName = '';
            state.organizationName = '';
        },
    },
});

// Action creators are generated for each case reducer function
export const {setProjectDetails, resetProjectDetails} = projectDetailsSlice.actions;

export default projectDetailsSlice.reducer;
