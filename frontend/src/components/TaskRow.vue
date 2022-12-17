<template>
  <div :class="['row', rowClasses]">
    <p>{{ $props.task.id }}</p>
    <p @click="toggleShowUpdate">{{ $props.task.name }}</p>
    <p>{{ $props.task.description }}</p>
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
      }
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

.row .check:hover, .row.done .check:hover {
  color: var(--fg-happy);
}

.row:nth-child(even) {
  background-color: var(--bg-color);
}

.row:nth-child(odd) {
  background-color: var(--bg-color-light);
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
</style>