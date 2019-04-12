<template>
  <div>
    <draggable v-model="parts" :move="movePart" @end="dropPart">
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
  },
  methods: {
    movePart(event) {
      this.$store.dispatch('setSortPart', {
        list: event.relatedContext.list,
        fromIndex: event.draggedContext.index,
        toIndex: event.draggedContext.futureIndex
      });
    },
    dropPart() {
      this.$eventBus.$emit('dropPart');
    }
  }
}
</script>

<style scoped>

</style>