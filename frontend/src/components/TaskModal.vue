<template>
  <Modal
      :title="title"
      :closeHandler="emitClosed"
  >
    <table class="task-modal" style="width: 100%">
      <tr>
        <td><label for="task-name">Name</label></td>
        <td><input id="task-name" type="text" v-model="requestState.name"/></td>
      </tr>
      <tr>
        <td><label for="task-desc">Description</label></td>
        <td><textarea id="task-desc" v-model="requestState.description"/></td>
      </tr>
      <tr>
        <td><label for="task-priority">Priority</label></td>
        <td>
          <select id="task-priority" v-model="requestState.priority">
            <option v-for="(name, idx) in priorityOptions" :value="5 - idx">{{name}}</option>
          </select>
          <span style="margin-left:16px;" v-text="priorityDescriptions[5 - requestState.priority]"></span>
        </td>
      </tr>
      <tr>
        <td></td>
        <td style="text-align: right">
          <button @click="submit" :disabled="submitDisabled">Update</button>
          <button class="secondary" @click="emitClosed">Cancel</button>
        </td>
      </tr>
    </table>
  </Modal>
</template>

<script>
import Modal from "./Modal.vue";

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
  name: "TaskModal",
  components: {Modal},
  emits: ['dialogClosed', 'taskUpdated', 'taskCreated'],
  props: {
    task: Object,
  },
  data() {
    return {
      requestState: {
        id: this.$props.task ? this.$props.task.id : 0,
        name: this.$props.task ? this.$props.task.name : '',
        description: this.$props.task ? this.$props.task.description : '',
        done: this.$props.task ? this.$props.task.done : false,
        priority: this.$props.task ? this.$props.task.priority : 0,
      }
    }
  },
  methods: {
    revertState() {
      this.requestState = {
        id: this.$props.task ? this.$props.task.id : 0,
        name: this.$props.task ? this.$props.task.name : '',
        description: this.$props.task ? this.$props.task.description : '',
        done: this.$props.task ? this.$props.task.done : false,
        priority: this.$props.task ? this.$props.task.priority : 0,
      };
    },
    emitClosed() {
      this.revertState();
      this.$emit('dialogClosed');
    },
    submit() {
      if (this.$props.task) {
        this.$emit('taskUpdated', this.requestState);
      } else {
        this.$emit('taskCreated', this.requestState);
        this.revertState();
      }
      this.$emit('dialogClosed');
    },
  },
  computed: {
    title() {
      if (this.$props.task) {
        return `Update task #${this.$props.task.id}`;
      }
      return 'Create Task';
    },
    priorityOptions() {
      return priorityOpts;
    },
    priorityDescriptions() {
      return priorityDesc
    },
    submitDisabled() {
      return this.requestState.name.length === 0;
    },
  }
}
</script>

<style scoped>
.task-modal {
  font-size: 1.2em;
}

.task-modal input[type="text"], .task-modal textarea {
  width: 100%;
  font-size: 1.2em;
}

.task-modal textarea {
  min-height: 150px;
}

.task-modal tr > td:first-child {
  width: 100px;
}
</style>