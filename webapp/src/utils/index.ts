/**
 * Utils
*/

import Constants from 'plugin_constants';

const getBaseUrls = (): {pluginApiBaseUrl: string; mattermostApiBaseUrl: string} => {
    const url = new URL(window.location.href);
    const baseUrl = `${url.protocol}//${url.host}`;
    const pluginUrl = `${baseUrl}/plugins/${Constants.pluginId}`;
    const pluginApiBaseUrl = `${pluginUrl}/api/v1`;
    const mattermostApiBaseUrl = `${baseUrl}/api/v4`;

    return {pluginApiBaseUrl, mattermostApiBaseUrl};
};

export const getCommandArgs = (command: string) => {
    const myRegexp = /[^\s"]+|"([^"]*)"/gi;
    const myArray = [];
    let match;
    do {
        match = myRegexp.exec(command);
        if (match != null) {
            myArray.push(match[1] ? match[1] : match[0]);
        }
    } while (match != null);
    return myArray.length > 2 ? myArray.slice(2) : [];
};

export const getProjectLinkModalArgs = (str: string): LinkPayload => {
    const data = str.split('/');
    if (data.length < 5 || (data[0] !== 'https:' && data[2] !== 'dev.azure.com')) {
        return {
            organization: '',
            project: '',
        };
    }

    return {
        organization: data[3] ?? '',
        project: data[4] ?? '',
    };
};

export const onPressingEnterKey = (event: Event | undefined, func: () => void) => {
    if (event instanceof KeyboardEvent && event.key !== 'Enter' && event.key !== ' ') {
        return;
    }

    func();
};

export const getProjectList = (data: ProjectDetails[]) => {
    const projectList: DropdownOptionType[] = [];
    data.map((project) => projectList.push({value: project.projectName, label: project.projectName}));
    return projectList;
};

export const getOrganizationList = (data: ProjectDetails[]) => {
    const uniqueOrganization: Record<string, boolean> = {};
    const organizationList: DropdownOptionType[] = [];
    data.map((organization) => {
        if (!(organization.organizationName in uniqueOrganization)) {
            uniqueOrganization[organization.organizationName] = true;
            organizationList.push({value: organization.organizationName, label: organization.organizationName});
        }
        return organizationList;
    });
    return organizationList;
};

export default {
    getBaseUrls,
};
