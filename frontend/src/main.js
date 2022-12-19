import {createApp} from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router';
import routes from './routes.js'
import App from './App.vue'
import confirm from './confirm.js'
import progress from "./progress.js";
import logEvent from "./logEvent.js";

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

confirm.setup();
progress.setup();
logEvent.setup();

let app = createApp(App);
app.use(router);
app.mount('#app');
