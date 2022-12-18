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
      <span class="material-icons delete" @click="$emit('taskDeleted', task.id)">delete</span>
    </div>
    <modal
        v-show="showUpdate"
        :title="updateTitle"
        @dialogClose="toggleShowUpdate"
    >
      <table class="modal-update" style="width: 100%">
        <tr>
          <td><label for="task-name">Name</label></td>
          <td><input id="task-name" type="text" v-model="updateState.name"/></td>
        </tr>
        <tr>
          <td><label for="task-desc">Description</label></td>
          <td><textarea id="task-desc" v-model="updateState.description"/></td>
        </tr>
        <tr>
          <td><label for="task-priority">Priority</label></td>
          <td>
            <select id="task-priority" v-model="updateState.priority">
              <option v-for="(name, idx) in priorityOptions" :value="5 - idx">{{name}}</option>
            </select>
            <span style="margin-left:16px;" v-text="priorityDescriptions[5 - updateState.priority]"></span>
          </td>
        </tr>
        <tr>
          <td></td>
          <td style="text-align: right">
            <button @click="submitUpdate" :disabled="updateDisabled">Update</button>
            <button class="secondary" @click="toggleShowUpdate">Cancel</button>
          </td>
        </tr>
      </table>
    </modal>
  </div>
</template>

<script>
import Modal from "./Modal.vue";

const descLimit = 50;
const nameLimit = 50;
const priorityOpts = ['Critical', 'High', 'Medium', 'Low', 'Lowest', 'None'];
const priorityDesc = [
  'Immediate priority, should be done ASAP',
  'Very important, should be done soon',
  'Moderately important and urgent',
  'Low importance and somewhat urgent',
  'Low importance and not urgent',
  'Routine, neither important nor urgent',
];

export default {
  name: "TaskRow",
  components: {Modal},
  emits: ["taskDone", "taskDeleted", "taskUpdate"],
  props: {
    task: Object,
  },
  data() {
    return {
      showDone: this.$props.task.done,
      showUpdate: false,
      updateState: {
        id: Number(this.$props.task.id),
        name: String(this.$props.task.name),
        description: String(this.$props.task.description),
        done: Boolean(this.$props.task.done),
        priority: Number(this.$props.task.priority),
      },
    }
  },
  methods: {
    taskDone() {
      this.showDone = !this.showDone;
      this.$emit('taskDone', this.$props.task.id);
    },
    revertUpdateState() {
      this.updateState = {
        id: Number(this.$props.task.id),
        name: String(this.$props.task.name),
        description: String(this.$props.task.description),
        done: Boolean(this.$props.task.done),
        priority: Number(this.$props.task.priority),
      }
    },
    toggleShowUpdate() {
      this.revertUpdateState();
      this.showUpdate = !this.showUpdate;
    },
    submitUpdate() {
      this.$emit("taskUpdate", this.updateState);
      this.toggleShowUpdate();
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
    updateTitle() {
      return "Update task #" + this.$props.task.id;
    },
    checkLigature() {
      if (this.showDone) {
        return "remove_done";
      }
      return "check";
    },
    updateDisabled() {
      return this.updateState.name === '';
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
    priorityOptions() {
      return priorityOpts;
    },
    priorityDescriptions() {
      return priorityDesc
    }
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

.row > *:nth-child(2):hover {
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

.row .material-icons {
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

.modal-update {
  font-size: 1.2em;
}

.modal-update input[type="text"], .modal-update textarea {
  width: 100%;
  font-size: 1.2em;
}

.modal-update textarea {
  min-height: 150px;
}

.modal-update tr > td:first-child {
  width: 100px;
}

div.name-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 1em;
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
  color: #2b2bff;
}

.row.done .material-icons.priority {
  color: gray;
}
</style>