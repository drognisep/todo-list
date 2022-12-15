<template>
  <loading v-if="isLoading"></loading>
  <div v-else class="content">
    <button id="new-task" @click.stop.prevent="this.newTasksShown = true">New Task</button>
    <div class="table">
      <div class="row header">
        <p>ID</p>
        <p>Task Name</p>
        <p>Description</p>
        <p>Status</p>
      </div>
      <div v-if="noTasks" class="row empty">
        <td colspan="4"><h3>No tasks found</h3></td>
      </div>
      <div v-else v-for="task in tasks" :key="task.id" class="row">
        <p>{{task.id}}</p>
        <p>{{ task.name }}</p>
        <p>{{ task.description }}</p>
        <p>{{ task.done ? "Done" : "Not Done" }}</p>
      </div>
    </div>
  </div>
  <modal v-show="newTasksShown" @dialogClose="toggleNewTasks" title="Create Task">
    <create-task @clickedCancel="toggleNewTasks" @taskCreated="addNewTask"></create-task>
  </modal>
</template>

<script>
import loadState from "../loadState.js";
import Loading from "../components/Loading.vue";
import {GetAllTasks} from "../wailsjs/go/main/TaskController.js";
import Modal from "../components/Modal.vue";
import CreateTask from "../components/CreateTask.vue";

export default {
  name: "AllTasks",
  components: {Loading, Modal, CreateTask},
  mixins: [loadState],
  data: () => {
    return {
      tasks: [],
      newTasksShown: false,
    }
  },
  methods: {
    getTasks() {
      this.startLoading();
      GetAllTasks()
          .then(tasks => {
            this.tasks = tasks.sort((a, b) => {
              if (a.id < b.id) {
                return -1;
              } else if (a.id > b.id) {
                return 1;
              }
              return 0;
            });
          })
          .catch(console.error)
          .then(() => {
            this.doneLoading();
          })
    },
    toggleNewTasks() {
      this.newTasksShown = !this.newTasksShown
    },
    addNewTask(data) {
      if(data) {
        if (!this.tasks) {
          this.tasks = [];
        }
        this.tasks.push(data);
      }
    }
  },
  computed: {
    noTasks() {
      return !this.tasks || this.tasks.length === 0
    },
  },
  created() {
    this.getTasks();
  }
}
</script>

<style scoped>
.content {
  --content-padding: 8px;

  padding: var(--content-padding);
  overflow-y: auto;
  max-height: 100%;
}

.table {
  width: 100%;
  display: table;
}

.table {
  --radius: 8px;
  --border-size: 2px;

  border-radius: var(--radius);
  border: var(--border-dark);
}

.table > .row:first-child > *:first-child {
  border-top-left-radius: calc(var(--radius) - 2px);
}

.table > .row:first-child > *:last-child {
  border-top-right-radius: calc(var(--radius) - 2px);
}

.table > .row:last-child > *:first-child {
  border-bottom-left-radius: calc(var(--radius) - 2px);
}

.table > .row:last-child > *:last-child {
  border-bottom-right-radius: calc(var(--radius) - 2px);
}

.row {
  display: table-row;
  border: 2px solid transparent;
}

.row > * {
  display: table-cell;
  padding: 4px 8px;
  background-color: transparent;
  border-right: var(--border-light);
}

.row > *:last-child {
  border-right: none;
}

.row > *:nth-child(1) {
  text-align: right;
}
.row > *:nth-child(2), .row > *:nth-child(3) {
  text-overflow: ellipsis;
  overflow: hidden;
}
.row > *:nth-child(4) {
  text-align: center;
}

.row:nth-child(even) {
  background-color: var(--bg-color);
}

.row:nth-child(odd) {
  background-color: var(--bg-color-light);
}

.row.header {
  font-weight: bold;
}

.row.header > * {
  text-align: center;
}

.row.empty > * {
  text-align: center;
}
</style>