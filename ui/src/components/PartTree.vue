<template>
  <div>
    <part-node v-for="part in parts" :key="part.ID" :option="part"></part-node>
  </div>
</template>

<script>
  import PartNode from "./PartNode";
  export default {
    name: "PartTree",
    components: {PartNode},
    props: {
      parts: {
        type: Array,
        required: true
      }
    },
    data() {
      return {
        selectedPart: {}
      }
    },
    methods: {
      selectPart(partId) {
        if (partId) {
          const partId = partId;
          this.$store.dispatch('setPartId', partId);
          this.$emit('select', partId);
        }
      },
    }
  }
</script>

<style scoped>

</style>