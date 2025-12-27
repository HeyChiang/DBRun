/// <reference types="vite/client" />

declare module '*.vue' {
    import type {DefineComponent} from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}

// Wails 全局类型声明
declare global {
    interface Window {
        runtime?: {
            WindowSetSize?: (width: number, height: number) => Promise<void>;
            [key: string]: any;
        };
        go?: {
            api?: {
                SQLiteAPI?: {
                    GetAllProjects?: () => Promise<any>;
                    [key: string]: any;
                };
                [key: string]: any;
            };
            [key: string]: any;
        };
    }
}
