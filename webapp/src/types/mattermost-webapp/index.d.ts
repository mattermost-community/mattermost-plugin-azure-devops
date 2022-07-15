/**
 * Keep all Mattermost webapp types here
 * To add more types you can refer https://developers.mattermost.com/extend/plugins/webapp/reference
 */

export interface PluginRegistry {
    registerPostTypeComponent(typeName: string, component: React.ElementType);
    registerReducer(reducer);
    registerRootComponent(component: ReactDOM);
    registerChannelIntroButtonAction(icon: JSX.Element, action: () => void, tooltipText?: string | null);
    registerChannelHeaderMenuAction(text: string, action: () => void);
    registerRightHandSidebarComponent(component: () => JSX.Element, title: string | JSX.Element);
    registerChannelHeaderButtonAction(icon: JSX.Element, action: () => void, dropdownText: string | null, tooltipText: string | null);
}
