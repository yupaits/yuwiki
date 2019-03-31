<template>
  <div>
    {{option.name}}
    <part-node v-for="part in option.SubParts"/>
  </div>
</template>

<script>
  export default {
    name: "PartNode",
    props: {
      option: {
        type: Object,
        required: true
      }
    }
  }
</script>

<style scoped>

</style>