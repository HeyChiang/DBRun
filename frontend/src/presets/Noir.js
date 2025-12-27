import { definePreset } from '@primevue/themes';
import Aura from '@primevue/themes/aura';

const Noir = definePreset(Aura, {
    semantic: {
        primary: {
            50: '{surface.50}',
            100: '{surface.100}',
            200: '{surface.200}',
            300: '{surface.300}',
            400: '{surface.400}',
            500: '{surface.500}',
            600: '{surface.600}',
            700: '{surface.700}',
            800: '{surface.800}',
            900: '{surface.900}',
            950: '{surface.950}'
        },
        icons: {
            light: {
                theme: 'pi pi-moon',
                themeLabel: '黑暗模式',
                file: 'pi pi-fw pi-file',
                new: 'pi pi-fw pi-plus',
                open: 'pi pi-fw pi-folder-open',
                save: 'pi pi-fw pi-save',
                edit: 'pi pi-fw pi-pencil',
                undo: 'pi pi-fw pi-undo',
                redo: 'pi pi-fw pi-refresh',
                view: 'pi pi-fw pi-eye',
                zoom: 'pi pi-fw pi-search',
                help: 'pi pi-fw pi-question-circle',
                about: 'pi pi-fw pi-info-circle'
            },
            dark: {
                theme: 'pi pi-sun',
                themeLabel: '明亮模式',
                file: 'pi pi-fw pi-file',
                new: 'pi pi-fw pi-plus',
                open: 'pi pi-fw pi-folder-open',
                save: 'pi pi-fw pi-save',
                edit: 'pi pi-fw pi-pencil',
                undo: 'pi pi-fw pi-undo',
                redo: 'pi pi-fw pi-refresh',
                view: 'pi pi-fw pi-eye',
                zoom: 'pi pi-fw pi-search',
                help: 'pi pi-fw pi-question-circle',
                about: 'pi pi-fw pi-info-circle'
            }
        },
        colorScheme: {
            light: {
                primary: {
                    color: '{primary.950}',
                    contrastColor: '#ffffff',
                    hoverColor: '{primary.900}',
                    activeColor: '{primary.800}'
                },
                highlight: {
                    background: '{primary.950}',
                    focusBackground: '{primary.700}',
                    color: '#ffffff',
                    focusColor: '#ffffff'
                }
            },
            dark: {
                primary: {
                    color: '{primary.50}',
                    contrastColor: '{primary.950}',
                    hoverColor: '{primary.100}',
                    activeColor: '{primary.200}'
                },
                highlight: {
                    background: '{primary.50}',
                    focusBackground: '{primary.300}',
                    color: '{primary.950}',
                    focusColor: '{primary.950}'
                }
            }
        }
    }
});

export default Noir;
