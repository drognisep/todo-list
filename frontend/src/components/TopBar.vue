<template>
<div id="topbar">
<h3 id="brand">Tasks</h3>
</div>
</template>

<script>
export default {
  name: "TopBar"
}
</script>

<style scoped>
#topbar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: var(--toolbar-height);
  background-color: var(--bg-lighter-color);
  box-shadow: black 0 0 5px;
}
#topbar > * {
  display: inline-block;
}
#topbar > #brand {
  font-weight: bold;
  height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  margin: 0 0 0 8px;
  cursor: default;
  text-shadow: black 0 0 8px;
}
</style>