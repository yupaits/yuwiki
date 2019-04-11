<template>
  <div>
    <draggable v-model="parts">
      <transition-group>
        <part-node v-for="part in parts" :key="part.ID" :option="part"></part-node>
      </transition-group>
    </draggable>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
import PartNode from "./PartNode";
export default {
  name: "PartTree",
  components: {draggable, PartNode},
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