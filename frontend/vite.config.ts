import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// svg
import svgLoader from 'vite-svg-loader'

// primevue
import Components from 'unplugin-vue-components/vite';
import {PrimeVueResolver} from '@primevue/auto-import-resolver';


// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
    plugins: [
        vue(),
        Components({
            resolvers: [
                PrimeVueResolver()
            ]
        }),
        svgLoader({
            svgoConfig: {
                multipass: true
            }
        })
    ],
    resolve: {
        alias: {
            '@': '/src',
        },
    },
    define: {
        '__VUE_PROD_HYDRATION_MISMATCH_DETAILS__': mode === 'production'
    }
})
)