export const getOrganizationList = (data: ProjectDetails[]) => {
    const uniqueOrganization = new Set();
    const organizationList: LabelValuePair[] = [];

    data.map((project) => uniqueOrganization.add(project.organizationName));

    uniqueOrganization.forEach((organization) => organizationList.push({
        value: organization as string,
        label: organization as string,
    }));

    return organizationList;
};

export const getProjectList = (data: ProjectDetails[]) => {
    const projectList: LabelValuePair[] = [];

    data.map((project) => projectList.push({value: project.projectName, label: project.projectName}));
    return projectList;
};

export const getCurrentChannelSubscriptions = (data: SubscriptionDetails[], channelID: string) => {
    const subscriptionList = (data || []).filter(((subscription) => subscription.channelID === channelID));
    return subscriptionList;
};

export const getCurrentChannelName = (data: ChannelList[], channelID: string) => {
    const currentChannel = (data || []).filter(((channel) => channel.id === channelID));
    return currentChannel[0].display_name;
};

