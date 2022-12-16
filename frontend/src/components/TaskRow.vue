<template>
  <div :class="rowClasses">
    <p>{{ $props.task.id }}</p>
    <p>{{ $props.task.name }}</p>
    <p>{{ $props.task.description }}</p>
    <div>
      <span class="material-icons check" @click="taskDone">check</span>
      <span class="material-icons delete" @click="$emit('taskDeleted', task.id)">delete</span>
    </div>
  </div>
</template>

<script>
export default {
  name: "TaskRow",
  emits: ["taskDone", "taskDeleted"],
  props: {
    task: Object,
  },
  data() {
    return {
      showDone: this.$props.task.done,
    }
  },
  methods: {
    taskDone() {
      this.showDone = !this.showDone;
      setTimeout(() => {
        this.$emit('taskDone', this.$props.task.id);
      }, 200);
    }
  },
  computed: {
    rowClasses() {
      if (this.showDone) {
        return "done row"
      }
      return "row"
    }
  }
}
</script>

<style scoped>
.row:last-child > *:first-child {
  border-bottom-left-radius: calc(var(--radius) - 2px);
}

.row:last-child > *:last-child {
  border-bottom-right-radius: calc(var(--radius) - 2px);
}

.row {
  display: table-row;
  border: 2px solid transparent;
}

.row > * {
  display: table-cell;
  padding: 4px 8px;
  background-color: transparent;
  border-right: var(--border-light);
}

.row > *:last-child {
  border-right: none;
}

.row > *:nth-child(1) {
  text-align: right;
}

.row > *:nth-child(2), .row > *:nth-child(3) {
  text-overflow: ellipsis;
  overflow: hidden;
}

.row > *:nth-child(4) {
  display: flex;
  position: relative;
  top: 3px;
  justify-content: center;
  text-align: center;
}

.row .material-icons {
  cursor: pointer;
}

.row .delete:hover, .row.done .delete:hover {
  color: var(--fg-danger);
}

.row .check:hover, .row.done .check:hover {
  color: var(--fg-happy);
}

.row:nth-child(even) {
  background-color: var(--bg-color);
}

.row:nth-child(odd) {
  background-color: var(--bg-color-light);
}

.row.done p {
  text-decoration: line-through;
  color: gray;
}

.row.done .material-icons {
  color: gray;
}
</style>