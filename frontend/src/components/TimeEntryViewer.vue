<template>
  <div class="entry-table">
    <div class="table-container" v-if="hasEntries && showTotals">
      <div>
        <h3 v-text="header"></h3>
        <table>
          <tr>
            <th>Day</th>
            <th>Start</th>
            <th>End</th>
            <th>Duration</th>
            <th>Task</th>
          </tr>
          <TimeEntryViewerRow v-for="entry in entries" :entry="entry"/>
        </table>
      </div>
      <div>
        <h3>{{header}} Summary</h3>
        <table>
          <tr>
            <th>Task</th>
            <th>Duration</th>
          </tr>
          <tr v-for="line in summaryLines" :key="line.name">
            <td>{{ summaryWrap(line.name) }}</td>
            <td>{{ line.duration }}</td>
          </tr>
        </table>
      </div>
    </div>
    <div v-else-if="hasEntries">
      <h3 v-text="header"></h3>
      <table>
        <tr>
          <th>Day</th>
          <th>Start</th>
          <th>End</th>
          <th>Duration</th>
          <th>Task</th>
        </tr>
        <TimeEntryViewerRow v-for="entry in entries" :entry="entry"/>
      </table>
    </div>
    <h3 v-else>No time entries found</h3>
  </div>
</template>

<script>
import TimeEntryViewerRow from "./TimeEntryViewerRow.vue";
import {GetSummaryForEntries} from "../wailsjs/go/main/TaskController.js";

export default {
  name: "TimeEntryViewer",
  components: {TimeEntryViewerRow},
  props: {
    header: String,
    entries: Array,
    showTotals: Boolean,
  },
  data() {
    return {
      summaryLines: [],
    };
  },
  methods: {
    summaryWrap(s) {
      if (s.length >= 30) {
        return s.substring(0, 27) + "...";
      }
      return s;
    },
  },
  computed: {
    hasEntries() {
      if (this.entries == null) {
        return false;
      }
      return this.entries.length > 0;
    },
  },
  watch: {
    entries: {
      handler(newEntries) {
        if (this.showTotals) {
          GetSummaryForEntries(newEntries)
              .then(summary => {
                this.summaryLines = summary.lines;
              })
              .catch(console.errorEvent)
        }
      },
      immediate: true,
    }
  },
}
</script>

<style scoped>
.entry-table {
  width: 100%;
}

.entry-table table {
  width: 100%;
}

th {
  font-weight: bold;
}

.table-container {
  display: flex;
  flex-direction: row;
}

.table-container > * {
  flex: 1;
}

.table-container > *:last-child {
  flex: 0.5;
}
</style>