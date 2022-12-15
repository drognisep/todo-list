<template>
  <table style="width:100%">
    <tbody>
    <tr>
      <td><label for="task-summary">Summary</label></td>
      <td><input id="task-summary" type="text" v-model="name"/></td>
    </tr>
    <tr>
      <td><label for="task-description">Description</label></td>
      <td><textarea id="task-description" v-model="description"></textarea></td>
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
      name: "",
      description: "",
      done: false,
    }
  },
  methods: {
    clickedSubmit() {
      let data = {
        name: this.name,
        description: this.description,
        done: this.done,
      }
      console.log("Creating task");
      console.log(data);
      CreateTask(data)
          .then(created => {
            console.log("Created task");
            console.log(created);
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
      this.name = "";
      this.description = "";
      this.done = false;
    },
  },
  computed: {
    submitDisabled() {
      return this.name === "";
    },
  }
}
</script>

<style scoped>
input[type="text"], textarea, input[type="text"]:focus, textarea:focus {
  background-color: var(--bg-color-light);
  color: var(--fg-color);
}
input[type="text"] {
  width: 100%;
}
textarea {
  width: 100%;
  min-height: 150px;
}
tr > td:first-child {
  width: 100px;
}
</style>