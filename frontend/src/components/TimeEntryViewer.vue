<template>
  <div class="entry-table">
    <h3 v-text="header"></h3>
    <table v-if="hasEntries">
      <tr>
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
      return this.entries.length > 0;
    }
  },
}
</script>

<style scoped>
th {
  font-weight: bold;
}
</style>