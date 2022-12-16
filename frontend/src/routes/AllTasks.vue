<template>
  <loading v-show="isLoading"></loading>
  <div class="content">
    <button v-show="!isLoading" id="new-task" @click.stop.prevent="this.newTasksShown = true">New Task</button>
    <div v-show="!isLoading" class="table">
      <div class="header">
        <p>ID</p>
        <p>Task Name</p>
        <p>Description</p>
        <p>Actions</p>
      </div>
      <div v-if="noTasks" class="empty">
        <td colspan="4"><h3>No tasks found</h3></td>
      </div>
      <TaskRow
          v-else
          v-for="task in tasks"
          :key="task.id"
          :task="task"
          @taskDone="toggleTaskDone"
          @taskDeleted="deleteTask"
          @taskUpdate="taskUpdate"
      />
    </div>
  </div>
  <modal v-show="newTasksShown" @dialogClose="toggleNewTasks" title="Create Task">
    <CreateTask @clickedCancel="toggleNewTasks" @taskCreated="taskCreated"></CreateTask>
  </modal>
</template>

<script>
import loadState from "../loadState.js";
import Loading from "../components/Loading.vue";
import {DeleteTask, GetAllTasks, UpdateTask} from "../wailsjs/go/main/TaskController.js";
import Modal from "../components/Modal.vue";
import CreateTask from "../components/CreateTask.vue";
import TaskRow from "../components/TaskRow.vue";

export default {
  name: "AllTasks",
  components: {Loading, Modal, CreateTask, TaskRow},
  mixins: [loadState],
  data() {
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
            this.sortTasks();
          })
          .catch(console.error)
          .then(() => {
            this.doneLoading();
          })
    },
    toggleNewTasks() {
      this.newTasksShown = !this.newTasksShown
    },
    taskCreated() {
      this.newTasksShown = false;
      this.getTasks();
    },
    toggleTaskDone(id) {
      this.tasks.forEach((item, i) => {
        if (item.id === id) {
          item.done = !item.done;
          this.startLoading();
          UpdateTask(id, item)
              .then(updated => {
                this.tasks[i] = updated;
                this.sortTasks();
              })
              .catch(console.error)
              .then(() => {
                this.doneLoading();
              })
        }
      })
    },
    taskUpdate(updated) {
      this.startLoading();
      console.log("Received update");
      console.log(updated);
      UpdateTask(updated.id, updated)
          .then(() => {
            this.getTasks();
          })
          .catch(console.error)
          .then(() => {
            this.doneLoading();
          })
    },
    deleteTask(id) {
      this.startLoading();
      DeleteTask(id)
          .catch(console.error)
          .then(() => {
            this.getTasks();
            this.doneLoading();
          })
    },
    sortTasks() {
      this.tasks = this.tasks.sort((a, b) => {
        if (a.done && !b.done) {
          return 1;
        } else if (!a.done && b.done) {
          return -1;
        }
        if (a.id < b.id) {
          return -1;
        } else if (a.id > b.id) {
          return 1;
        }
        return 0;
      });
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

.header {
  font-weight: bold;
  background-color: var(--bg-color-light);
}

.header > * {
  text-align: center;
}

.header > *:first-child {
  border-top-left-radius: calc(var(--radius) - 2px);
}

.header > *:last-child {
  border-top-right-radius: calc(var(--radius) - 2px);
}

.empty > * {
  text-align: center;
}

.empty, .header {
  display: table-row;
  border: 2px solid transparent;
}

.empty > *, .header > * {
  display: table-cell;
  padding: 4px 8px;
  background-color: transparent;
  border-right: var(--border-light);
}
</style>