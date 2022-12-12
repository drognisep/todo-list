import Placeholder from "./routes/Placeholder.vue";
import Dashboard from "./routes/Dashboard.vue";

export default [
    {path: '/', name: 'dashboard', component: Dashboard},
    {path: '/allTasks', name: 'allTasks', component: Placeholder},
    {path: '/searchTasks', name: 'searchTasks', component: Placeholder},
];
