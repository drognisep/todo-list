<template>
  <loading v-if="waiting !== 0"/>
  <div v-if="waiting === 0" class="container">
    <div class="line"><h3>Task Count: </h3>
      <p>{{ count }}</p></div>
  </div>
</template>

<script>
import {Count} from "../wailsjs/go/main/TaskController.js";
import Loading from "../components/Loading.vue";

export default {
  name: "Dashboard",
  components: {Loading},
  data: () => {
    return {
      waiting: 0,
      count: 0,
    }
  },
  created() {
    this.taskCount();
  },
  methods: {
    loading() {
      this.waiting++;
    },
    doneLoading() {
      if (this.waiting > 0) {
        this.waiting--;
      }
    },
    taskCount() {
      this.loading()
      Count()
          .then(count => {
            this.count = count;
          })
          .catch(err => {
            console.error("Error loading task count: " + err);
          })
          .then(() => this.doneLoading());
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