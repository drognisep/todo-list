<template>
  <Modal
      title="Import"
      :closeHandler="emitClosed"
  >
    <p>Use the drop-down box below to select what will happen if there are merge conflicts between your store and the import file.</p>
    <p>Click the <b>Import</b> button to select a file and start the import. Click the <b>Cancel</b> button to close this dialog and do nothing.</p>
    <table style="width: 100%">
      <tr>
        <td style="width: 100px"><label for="merge-strategy">Merge Strategy</label></td>
        <td>
          <select style="width: 100%" id="merge-strategy" v-model="requestState.strategy">
            <option v-for="(strat, idx) in mergeStrategies" :value="strat">{{ mergeDesc[idx] }}</option>
          </select>
        </td>
      </tr>
      <tr>
        <td></td>
        <td style="text-align: right">
          <button @click="submit">Import</button>
          <button class="secondary" @click="emitClosed">Cancel</button>
        </td>
      </tr>
    </table>
  </Modal>
</template>

<script>
import Modal from "./Modal.vue";

export default {
  name: "ImportDialog",
  components: {Modal},
  emits: ['dialogClosed','doImport'],
  data() {
    return {
      requestState: {
        strategy: 'KeepInternal',
      },
    };
  },
  methods: {
    emitClosed() {
      this.revertState();
      this.$emit('dialogClosed');
    },
    revertState() {
      this.requestState = {
        strategy: 'KeepInternal',
      };
    },
    submit() {
      this.$emit('doImport', this.requestState);
      this.emitClosed();
    },
  },
  computed: {
    mergeStrategies() {
      return ['KeepInternal', 'Overwrite', 'Error', 'Append'];
    },
    mergeDesc() {
      return [
        "Do not overwrite the state in your task store",
        "Overwrite data in your task store if there's a conflict",
        "Return an error if there's a conflict",
        "Append data, use this if the snapshot is from someone else",
      ];
    }
  },
}
</script>

<style scoped>

</style>