<template>
  <div id="topbar">
    <h3 id="brand-text">Tasks</h3>
    <p class="link">All Tasks</p>
    <p class="link">Task Search</p>
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
  --text-shadow: black 0 0 8px
}

#topbar > * {
  display: inline-block;
  max-height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  user-select: none;
  margin: 0;
  padding: 0;
}

#topbar > #brand-text {
  font-weight: bold;
  height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  margin: 0 32px 0 16px;
  cursor: default;
  text-shadow: var(--text-shadow);
}

#topbar > .link {
  position: relative;
  top: -1px;
  padding: 0 8px;
  margin: 0;
  cursor: pointer;
  transition: all 200ms linear;
}

#topbar > .link:hover, #topbar > .link.active {
  text-shadow: var(--text-shadow);
  color: var(--fg-highlight);
  background-color: var(--bg-highlight);
}
</style>