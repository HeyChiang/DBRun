import {createApp} from 'vue'
import App from './App.vue'
import router from './router'

import './style.css';
import "./flags.css";
import "primeicons/primeicons.css";
import PrimeVue from 'primevue/config';
import ConfirmationService from 'primevue/confirmationservice'
import DialogService from 'primevue/dialogservice'
import ToastService from 'primevue/toastservice';

import { createPinia } from 'pinia'

import appState from "./plugins/appState.js";
import Noir from './presets/Noir.js';

const app = createApp(App);

const pinia = createPinia()
app.use(pinia)
app.use(router)

app.use(appState)
app.use(ConfirmationService);
app.use(ToastService);
app.use(DialogService);
app.use(PrimeVue, {
    theme: {
        preset: Noir,
        options: {
            prefix: 'p',
            darkModeSelector: 'p-dark',
            cssLayer: false
        }
    }
});

app.mount('#app')
