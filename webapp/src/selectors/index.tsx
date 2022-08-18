export const getprojectDetailsState = (state: ReduxState) => {
    return state['plugins-mattermost-plugin-azure-devops'].projectDetailsSlice;
};

export const getRhsState = (state: ReduxState): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};
