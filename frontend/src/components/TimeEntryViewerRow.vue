<template>
  <tr>
    <td><p>{{ dayOfWeek }}</p></td>
    <td><p>{{ startTime }}</p></td>
    <td><p>{{ endTime }}</p></td>
    <td><p>{{ duration }}</p></td>
    <td><p>{{ nameWrap(taskName) }}</p></td>
  </tr>
</template>

<script>
import {GetTaskByID} from "../wailsjs/go/main/ModelController.js";
import {durationGo, formatClockTime, weekdaySemantic} from "../datetime.js";
import loadState from "../loadState.js";

export default {
  name: "TimeEntryViewerRow",
  mixins: [loadState],
  props: {
    entry: Object,
  },
  data() {
    return {
      entryState: null,
      taskState: null,
    };
  },
  methods: {
    nameWrap(s) {
      if (s.length >= 30) {
        return s.substring(0, 27) + "...";
      }
      return s;
    },
  },
  computed: {
    dayOfWeek() {
      if (this.entryState == null) {
        return '';
      }
      return weekdaySemantic(this.entryState.start);
    },
    startTime() {
      if (this.entryState == null) {
        return '';
      }
      return formatClockTime(this.entryState.start);
    },
    endTime() {
      if (this.entryState == null) {
        return '';
      }
      if (this.entryState.end == null) {
        return 'Tracking...';
      }
      return formatClockTime(this.entryState.end);
    },
    duration() {
      if (this.entryState == null) {
        return '';
      }
      if (this.entryState.end == null) {
        return '';
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
        this.startLoading();
        GetTaskByID(newEntry.taskID)
            .then(task => {
              this.taskState = task;
              this.entryState = newEntry;
            })
            .catch(console.logEvent)
            .then(() => {
              this.doneLoading();
            });
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

td p {
  user-select: text;
}
</style>