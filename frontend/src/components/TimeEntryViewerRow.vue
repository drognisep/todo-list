<template>
  <tr>
    <td><p>{{ duration }}</p></td>
    <td><p>{{ taskName }}</p></td>
  </tr>
</template>

<script>
import {GetTaskByID} from "../wailsjs/go/main/TaskController.js";
import {durationGo} from "../datetime.js";

export default {
  name: "TimeEntryViewerRow",
  props: {
    entry: Object,
  },
  data() {
    return {
      entryState: null,
      taskState: null,
    };
  },
  computed: {
    duration() {
      if (this.entryState == null) {
        return '';
      }
      if (this.entryState.end == null) {
        return 'Tracking...';
      }
      return durationGo(this.entryState.start, this.entryState.end);
    },
    taskName() {
      if (this.taskState == null) {
        return '';
      }
      return this.taskState.name;
    },
  },
  watch: {
    entry: {
      handler(newEntry) {
        if (newEntry == null) {
          this.taskState = null;
          this.entryState = null;
        }
        GetTaskByID(newEntry.taskID)
            .then(task => {
              this.taskState = task;
              this.entryState = newEntry;
            })
            .catch(console.logEvent);
      },
      immediate: true,
    }
  }
}
</script>

<style scoped>
td > * {
  margin: 0 0 0 16px;
}
td:nth-child(1) > * {
  margin-left: 0;
}
</style>