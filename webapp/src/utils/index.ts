/**
 * Utils
*/
import Constants from 'pluginConstants';

import {getOrganizationList, getProjectList} from './filterData';
import getErrorMessage from './errorHandling';

const getBaseUrls = (): { pluginApiBaseUrl: string; mattermostApiBaseUrl: string } => {
    const url = new URL(window.location.href);
    const baseUrl = `${url.protocol}//${url.host}`;
    const pluginUrl = `${baseUrl}/plugins/${Constants.common.pluginId}`;
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
        organization: decodeURI(data[3]) ?? '',
        project: decodeURI(data[4]) ?? '',
    };
};

export const getCreateTaskModalCommandArgs = (arr: Array<string>): CreateTaskCommandArgs => ({
    title: arr[1] ?? '',
    description: arr[2] ?? '',
});

export const onPressingEnterKey = (event: React.KeyboardEvent<HTMLSpanElement> | React.KeyboardEvent<SVGSVGElement>, func: () => void) => {
    if (event.key !== 'Enter' && event.key !== ' ') {
        return;
    }

    func();
};

export const sortProjectList = (project1: ProjectDetails, project2: ProjectDetails) => project1.projectName.toLocaleLowerCase().localeCompare(project2.projectName.toLocaleLowerCase());

export const addPathParamsToApiUrl = (url: string, pathParams?: Record<string, string>) => {
    if (!pathParams) {
        return url;
    }

    let newUrl = url;
    Object.keys(pathParams).forEach((param) => {
        newUrl = newUrl.replace(`:${param}`, pathParams[param]);
    });
    return newUrl;
};

export const formLabelValuePair = (labelKey: string, valueKey: string, data: Record<string, string>) => {
    const labelValuePair: LabelValuePair = {
        label: data[labelKey] ?? '',
        value: data[valueKey] ?? '',
    };

    return labelValuePair;
};

export const formLabelValuePairs = (labelKey: string, valueKey: string, data: Record<string, string>[]) => {
    const labelValuePairs: LabelValuePair[] = [];
    data.forEach((item) => labelValuePairs.push(formLabelValuePair(labelKey, valueKey, item)));

    return labelValuePairs;
};

export default {
    getBaseUrls,
    getErrorMessage,
    getOrganizationList,
    getProjectList,
    addPathParamsToApiUrl,
    formLabelValuePair,
    formLabelValuePairs,
};
