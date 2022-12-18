<template>
  <div :class="['row', rowClasses]">
    <p>{{ $props.task.id }}</p>
    <div class="name-cell">
      <p @click="toggleShowUpdate">{{ truncatedName }}</p>
      <span :class="['material-icons', 'priority', priorityClass]">{{ priorityIcon }}</span>
    </div>
    <p>{{ truncatedDesc }}</p>
    <div class="actions">
      <span class="material-icons check" @click="taskDone">{{ checkLigature }}</span>
      <span :class="['material-icons','favorite',$props.task.favorite ? 'favored' : '']" @click="taskFavorite">push_pin</span>
      <span class="material-icons delete" @click="deleteTask">delete</span>
    </div>
    <TaskModal
        v-show="showUpdate"
        :task="$props.task"
        @dialogClosed="toggleShowUpdate"
        @taskUpdated="submitUpdate"
    />
  </div>
</template>

<script>
import Modal from "./Modal.vue";
import TaskModal from "./TaskModal.vue";

const descLimit = 50;
const nameLimit = 50;

export default {
  name: "TaskRow",
  components: {Modal, TaskModal},
  emits: ["taskDeleted", "taskUpdate"],
  props: {
    task: Object,
  },
  data() {
    return {
      showDone: this.$props.task.done,
      showUpdate: false,
    }
  },
  methods: {
    taskDone() {
      this.showDone = !this.showDone;
      let data = this.$props.task;
      data.favorite = false;
      data.done = this.showDone;
      this.submitUpdate(data);
    },
    taskFavorite() {
      let data = this.$props.task;
      data.favorite = !data.favorite;
      this.submitUpdate(data);
    },
    toggleShowUpdate() {
      this.showUpdate = !this.showUpdate;
    },
    submitUpdate(state) {
      this.$emit("taskUpdate", state);
    },
    deleteTask() {
      if (!this.$props.task.done) {
        confirmDialog('Are you sure?', 'This task is not completed and will be permanently deleted. Click OK to continue.', () => {
          this.$emit('taskDeleted', this.$props.task.id);
        })
        return;
      }
      this.$emit('taskDeleted', this.$props.task.id);
    }
  },
  computed: {
    rowClasses() {
      if (this.showDone) {
        return "done"
      }
      return ""
    },
    truncatedName() {
      if (this.$props.task.name.length > nameLimit) {
        return this.$props.task.name.substring(0, nameLimit) + "...";
      }
      return this.$props.task.name;
    },
    truncatedDesc() {
      if (this.$props.task.description.length > descLimit) {
        return this.$props.task.description.substring(0, descLimit) + "...";
      }
      return this.$props.task.description;
    },
    checkLigature() {
      if (this.showDone) {
        return "remove_done";
      }
      return "check";
    },
    priorityClass() {
      switch (this.$props.task.priority) {
        case 5:
          return "priority-highest";
        case 4:
          return "priority-high";
        case 3:
          return "priority-med";
        case 2:
          return "priority-low";
        case 1:
          return "priority-lowest";
        case 0:
        default:
          return "priority-none";
      }
    },
    priorityIcon() {
      switch (this.$props.task.priority) {
        case 5:
          return "keyboard_double_arrow_up";
        case 4:
          return "keyboard_arrow_up";
        case 3:
          return "unfold_more";
        case 2:
          return "keyboard_arrow_down";
        case 1:
          return "keyboard_double_arrow_down";
        case 0:
        default:
          return "";
      }
    },
  }
}
</script>

<style scoped>
.row:last-child > *:first-child {
  border-bottom-left-radius: calc(var(--radius) - 2px);
}

.row:last-child > *:last-child {
  border-bottom-right-radius: calc(var(--radius) - 2px);
}

.row {
  display: table-row;
  border: 2px solid transparent;
}

.row > * {
  display: table-cell;
  padding: 4px 8px;
  border-right: var(--border-light);
}

.row > *:nth-child(1) {
  text-align: right;
}

.row > *:nth-child(2) p:hover {
  color: var(--fg-highlight);
  cursor: pointer;
  text-decoration: underline;
}

.row > *:nth-child(2), .row > *:nth-child(3) {
  text-overflow: ellipsis;
  overflow: hidden;
}

.row > *:nth-child(4) {
  display: flex;
  position: relative;
  top: 3px;
  justify-content: center;
  text-align: center;
  border-right: none;
}

.row .material-icons:not(*.priority) {
  cursor: pointer;
}

.row .delete:hover, .row.done .delete:hover {
  color: var(--fg-danger);
}

.row .check:hover {
  color: var(--fg-happy);
}

.row.done .check:hover {
  color: var(--fg-warn);
}

.row .favorite.favored {
  color: var(--fg-blue);
}

.row .favorite:hover {
  color: var(--fg-blue);
}

.row:nth-child(even) {
  background-color: var(--bg-color);
}

.row:nth-child(odd) {
  background-color: var(--bg-color-light);
}

.row:nth-child(even):hover, .row:nth-child(odd):hover {
  background-color: var(--bg-highlight);
}

.row.done p {
  text-decoration: line-through;
  color: gray;
}

.row.done .material-icons {
  color: gray;
}

div.name-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 1em;
}

.material-icons {
  text-shadow: 0 0 3px black;
}

.material-icons.priority.priority-highest {
  color: var(--fg-danger);
}

.material-icons.priority.priority-high {
  color: var(--fg-warn);
}

.material-icons.priority.priority-med {
  color: yellow;
}

.material-icons.priority.priority-low {
  color: var(--fg-happy);
}

.material-icons.priority.priority-lowest {
  color: var(--fg-blue);
}

.row.done .material-icons.priority {
  color: gray;
}
</style>