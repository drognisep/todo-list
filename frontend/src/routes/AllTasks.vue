<template>
    <loading v-if="isLoading"></loading>
    <div v-if="!isLoading" class="content">
      <button id="new-task" @click.stop.prevent="this.newTasksShown = true">New Task</button>
      <div class="table">
        <div class="row header">
          <p>Task Name</p>
          <p>Description</p>
          <p>Status</p>
        </div>
        <div v-if="noTasks" class="row empty">
          <td colspan="3"><h3>No tasks found</h3></td>
        </div>
        <div v-else v-for="task in tasks" class="row">
          <p>{{ task.name }}</p>
          <p>{{ task.description }}</p>
          <p style="text-align: center">{{ task.done ? "Done" : "Not Done" }}</p>
        </div>
      </div>
    </div>
    <modal v-show="newTasksShown" @dialogClose="toggleNewTasks" title="Test">
      <button>Do a thing!</button>
    </modal>
</template>

<script>
import loadState from "../loadState.js";
import Loading from "../components/Loading.vue";
import {GetAllTasks} from "../wailsjs/go/main/TaskController.js";
import Modal from "../components/Modal.vue";

export default {
  name: "AllTasks",
  components: {Loading, Modal},
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
            this.tasks = tasks;
          })
          .catch(console.error)
          .then(() => {
            this.doneLoading();
          })
    },
    toggleNewTasks() {
      this.newTasksShown = !this.newTasksShown
    },
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