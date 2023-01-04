<template>
  <loading v-if="isLoading"/>
  <div v-else class="container">
    <div class="line"><p>Task Count: </p>
      <p>{{ count }}</p>
    </div>
    <TimeEntryViewer header="Today's Time Entries" :entries="todayEntries"/>
  </div>
</template>

<script>
import {Count, GetTimeEntriesToday} from "../wailsjs/go/main/TaskController.js";
import Loading from "../components/Loading.vue";
import loading from "../loadState.js";
import TimeEntryViewer from "../components/TimeEntryViewer.vue";

export default {
  name: "Dashboard",
  components: {TimeEntryViewer, Loading},
  mixins: [loading],
  data() {
    return {
      waiting: 0,
      count: 0,
      day: [],
    }
  },
  methods: {
    taskCount() {
      this.startLoading()
      Count()
          .then(count => {
            this.count = count;
          })
          .catch(err => {
            console.errorEvent("Error loading task count: " + err);
          })
          .then(() => this.doneLoading());
    },
    timeEntryDetails() {
      this.startLoading();
      GetTimeEntriesToday()
          .then(entries => {
            this.day = entries;
          })
          .catch(console.errorEvent)
          .then(this.doneLoading);
    }
  },
  computed: {
    todayEntries() {
      return this.day;
    }
  },
  created() {
    this.taskCount();
    this.timeEntryDetails();
  },
}
</script>

<style scoped>
.container {
  padding: 8px;
  overflow-y: auto;
  max-height: 100%;
}

.line > * {
  display: inline-block;
  padding-left: 8px;
}

.line > *:first-child {
  padding-left: 0;
}
</style>