<template>
  <div id="topbar">
    <img alt="tasks" id="brand-img" src="../assets/images/gear.svg"/>
    <h3 id="brand-text">Tasks</h3>
    <RouterLink to="/" class="link">Dashboard</RouterLink>
    <RouterLink to="/allTasks" class="link">All Tasks</RouterLink>
    <RouterLink to="/searchTasks" class="link">Task Search</RouterLink>
  </div>
</template>

<script>
import {RouterLink} from 'vue-router';

export default {
  name: "TopBar",
  components: [RouterLink],
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
  padding: 0 0 0 16px;
  z-index: var(--z-toolbar);

  --text-shadow: black 0 0 8px;
  --border-right: 2px solid var(--overlay-color);
}

#topbar > * {
  display: inline-block;
  max-height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  margin: 0;
  padding: 0;
}

#topbar > #brand-img {
  height: calc(var(--toolbar-height) - 12px);
  margin: 6px 16px 0 0;
  float: left;
}

#topbar > #brand-text {
  font-weight: bold;
  height: var(--toolbar-height);
  line-height: var(--toolbar-height);
  padding: 0 16px 0 0;
  cursor: default;
  text-shadow: var(--text-shadow);
  border-right: var(--border-right);
}

#topbar > .link {
  color: var(--fg-color);
  text-decoration: none;
  position: relative;
  top: -1px;
  padding: 0 8px;
  margin: 0;
  cursor: pointer;
  transition: all 200ms linear;
  border-right: var(--border-right);
}

#topbar > .link:hover, #topbar > .link.router-link-active {
  text-shadow: var(--text-shadow);
  color: var(--fg-highlight);
  background-color: var(--bg-highlight);
}
</style>