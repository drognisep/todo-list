<template>
  <div class="entry-table">
    <h3 v-text="header"></h3>
    <table v-if="hasEntries">
      <tr>
        <th>Day</th>
        <th>Start</th>
        <th>End</th>
        <th>Duration</th>
        <th>Task</th>
      </tr>
      <TimeEntryViewerRow v-for="entry in entries" :entry="entry" />
    </table>
    <h3 v-else>No time entries found</h3>
  </div>
</template>

<script>
import TimeEntryViewerRow from "./TimeEntryViewerRow.vue";

export default {
  name: "TimeEntryViewer",
  components: {TimeEntryViewerRow},
  props: {
    header: String,
    entries: Array,
  },
  methods: {
  },
  computed: {
    hasEntries() {
      if (this.entries == null) {
        return false;
      }
      return this.entries.length > 0;
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
</style>