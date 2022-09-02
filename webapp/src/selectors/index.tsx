export const getprojectDetailsState = (state: ReduxState) => {
    return state['plugins-mattermost-plugin-azure-devops'].projectDetailsSlice;
};

export const getRhsState = (state: ReduxState): {isSidebarOpen: boolean} => {
    return state.views.rhs;
};

export const getGlobalModalState = (state: ReduxState): GlobalModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].globalModalSlice;
};

export const getLinkModalState = (state: ReduxState): LinkProjectModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].openLinkModalSlice;
};

export const getTaskModalState = (state: ReduxState): CreateTaskModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].openTaskModalSlice;
};

export const getSubscribeModalState = (state: ReduxState): SubscribeModalState => {
    return state['plugins-mattermost-plugin-azure-devops'].openSubscribeModalSlice;
};
