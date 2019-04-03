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
    }
  }
</script>

<style scoped>

</style>