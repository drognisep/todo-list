import Dashboard from "./routes/Dashboard.vue";
import AllTasks from "./routes/AllTasks.vue";

export default [
    {path: '/', name: 'dashboard', component: Dashboard},
    {path: '/allTasks', name: 'allTasks', component: AllTasks},
];
