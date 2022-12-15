<template>
  <loading v-if="loading"/>
  <div v-else class="container">
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
    startLoading() {
      this.waiting++;
    },
    doneLoading() {
      if (this.waiting > 0) {
        this.waiting--;
      }
    },
    taskCount() {
      this.startLoading()
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
  computed: {
    loading() {
      return this.waiting !== 0;
    },
  }
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