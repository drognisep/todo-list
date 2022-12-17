<template>
  <table class="modal-create">
    <tbody>
    <tr>
      <td><label for="task-summary">Summary</label></td>
      <td><input id="task-summary" type="text" v-model="createState.name"/></td>
    </tr>
    <tr>
      <td><label for="task-description">Description</label></td>
      <td><textarea id="task-description" v-model="createState.description"></textarea></td>
    </tr>
    <tr>
      <td colspan="2">
        <div style="text-align: right">
          <button :disabled="submitDisabled" @click.prevent.stop="clickedSubmit">Create</button>
          <button class="secondary" @click.prevent.stop="clickedCancel">Close</button>
        </div>
      </td>
    </tr>
    </tbody>
  </table>
</template>

<script>
import {CreateTask} from "../wailsjs/go/main/TaskController.js";

export default {
  name: "CreateTask",
  data() {
    return {
      createState: {
        name: "",
        description: "",
        done: false,
      },
    }
  },
  methods: {
    clickedSubmit() {
      CreateTask(this.createState)
          .then(created => {
            this.$emit("taskCreated", created);
          })
          .then(() => {
            this.resetState()
          })
          .catch(console.error)
    },
    clickedCancel() {
      this.resetState();
      this.$emit("clickedCancel");
    },
    resetState() {
      this.createState = {
        name: "",
        description: "",
        done: false,
      }
    },
  },
  computed: {
    submitDisabled() {
      return this.createState.name === "";
    },
  }
}
</script>

<style scoped>
.modal-create {
  width: 100%;
  font-size: 1.2em;
}

.modal-create input[type="text"], .modal-create textarea {
  font-size: 1.2em;
}

.modal-create textarea {
  min-height: 150px;
}

tr > td:first-child {
  width: 100px;
}
</style>