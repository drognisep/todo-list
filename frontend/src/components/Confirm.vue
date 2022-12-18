<template>
  <div v-show="showConfirm" class="outer" @click.stop.prevent>
    <div class="inner">
      <div class="confirm">
        <h1>{{ title }}</h1>
        <p>{{ message }}</p>
        <div class="btnDiv">
          <button @click="confirm">OK</button>
          <button @click="hideConfirm">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Confirm",
  data() {
    return {
      title: '',
      message: '',
      onConfirm: () => {},
      showConfirm: false,
      listener: null,
    }
  },
  methods: {
    confirm() {
      this.onConfirm();
      this.hideConfirm();
    },
    hideConfirm() {
      this.showConfirm = false;
    },
    showConfirmHandler(evt) {
      this.title = evt.detail.title;
      this.message = evt.detail.message;
      this.onConfirm = evt.detail.onConfirm;
      this.showConfirm = true;
    }
  },
  created() {
    window.addEventListener('confirmDialog', (evt) => {
      this.showConfirmHandler(evt);
    });
  },
  destroyed() {
    window.removeEventListener('confirmDialog', this.showConfirmHandler.bind(this));
  }
}
</script>

<style scoped>
.outer {
  position: fixed;
  z-index: var(--z-toolbar);
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--overlay-color);
}

.inner {
  position: absolute;
  left: 25vw;
  width: 50vw;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.confirm {
  display: inline-block;
  background-color: var(--bg-color);
  border-radius: 8px;
  padding: 16px;
}

.btnDiv {
  text-align: right;
}
</style>