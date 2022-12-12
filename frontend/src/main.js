import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router';
import routes from './routes.js'
import App from './App.vue'

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

let app = createApp(App);
app.use(router);
app.mount('#app');
