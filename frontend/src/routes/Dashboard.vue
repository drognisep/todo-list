<template>
  <div class="container">
    <div class="line"><h3>Task Count: </h3>
      <p>{{ count }}</p></div>
  </div>
</template>

<script>
import { Count } from "../wailsjs/go/main/TaskController.js";

export default {
  name: "Dashboard",
  data: () => {
    return {
      count: -1,
    }
  },
  created() {
    this.taskCount()
  },
  methods: {
    taskCount() {
      return Count().then(count => {
        this.count = count;
      }).catch(err => console.error);
    }
  },
}
</script>

<style scoped>
.container {
  padding: 8px;
}

.line > * {
  display: inline-block;
  padding-left: 8px;
}

.line > *:first-child {
  padding-left: 0;
}
</style>