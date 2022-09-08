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

    // TODO: add proper structure for Organization in metaData field
    data.map((project) => projectList.push({value: project.projectName, label: project.projectName, metaData: project.organizationName}));
    return projectList;
};

export const getCurrentChannelSubscriptions = (data: SubscriptionDetails[], channelID: string): SubscriptionDetails[] => data.filter(((subscription) => subscription.channelID === channelID));
