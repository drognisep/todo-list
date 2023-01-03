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
          :tracking="task.id === trackedTaskID"
          @taskDeleted="deleteTask"
          @taskUpdate="taskUpdate"
      />
    </div>
  </div>
  <TaskModal
      v-show="newTasksShown"
      @dialogClosed="toggleNewTasks"
      @taskCreated="taskCreated"
  />
</template>

<script>
import loadState from "../loadState.js";
import Loading from "../components/Loading.vue";
import {
  CreateTask,
  DeleteTask,
  GetAllTasks,
  GetTrackedTaskDetails,
  UpdateTask
} from "../wailsjs/go/main/TaskController.js";
import Modal from "../components/Modal.vue";
import TaskRow from "../components/TaskRow.vue";
import TaskModal from "../components/TaskModal.vue";
import {EventsOn} from "../wailsjs/runtime/runtime.js";

export default {
  name: "AllTasks",
  components: {Loading, Modal, TaskModal, TaskRow},
  mixins: [loadState],
  data() {
    return {
      tasks: [],
      newTasksShown: false,
      trackedTaskID: 0,
    }
  },
  methods: {
    getTasks() {
      GetAllTasks()
          .then(tasks => {
            this.tasks = tasks;
          })
          .catch(console.errorEvent)
    },
    toggleNewTasks() {
      this.newTasksShown = !this.newTasksShown
    },
    taskCreated(newTask) {
      CreateTask(newTask)
          .then(() => {
            this.getTasks();
          })
          .catch(console.errorEvent)
    },
    taskUpdate(updated) {
      UpdateTask(updated.id, updated)
          .then(() => {
            this.getTasks();
          })
          .catch(console.errorEvent)
    },
    deleteTask(id) {
      DeleteTask(id)
          .catch(console.errorEvent)
          .then(() => {
            this.getTasks();
          })
    },
    onTaskStarted(evt) {
      this.trackedTaskID = evt.started;
    },
    onTaskStopped() {
      this.trackedTaskID = 0;
    },
  },
  computed: {
    noTasks() {
      return !this.tasks || this.tasks.length === 0
    },
  },
  created() {
    this.startLoading();
    this.startLoading();
    GetAllTasks()
        .then(tasks => {
          this.tasks = tasks;
        })
        .catch(console.errorEvent)
        .then(() => {
          this.doneLoading();
        });
    GetTrackedTaskDetails()
        .then(details => {
          this.trackedTaskID = details.task.id;
        })
        .catch(console.errorEvent)
        .then(() => {
          this.doneLoading();
        });
    EventsOn("taskStarted", this.onTaskStarted);
    EventsOn("taskStopped", this.onTaskStopped);
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