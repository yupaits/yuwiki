<template>
  <draggable v-model="parts" :move="movePart" @end="dropPart">
    <transition-group>
      <div v-for="part in parts" :key="part.ID">
        <div class="part-item" :class="{'active': $store.getters.partId === part.ID}" @click="selectPart(part)">
          <a-icon :type="part.partType === 0 ? 'folder-open' : 'inbox'"/> {{part.name}}
          <span class="pull-right" v-if="part.star"><a-icon type="star" theme="filled" :style="{color: '#fadb14'}"/></span>
        </div>
        <part-tree :parts="part.subParts" class="ml-2"></part-tree>
      </div>
    </transition-group>
  </draggable>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  name: "PartTree",
  components: {draggable},
  props: {
    parts: {
      type: Array,
      required: true
    }
  },
  methods: {
    selectPart(part) {
      const partId = part.ID;
      this.$store.dispatch('setPartId', partId);
      this.$store.dispatch('setPartStar', part.star);
      delete part.SubParts;
      this.$store.dispatch('setPart', part);
      this.$emit('select', partId);
      this.$eventBus.$emit('selectPart', partId);
    },
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
.part-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.part-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.part-item.active {
  background: #91d5ff;
  font-weight: bold;
  color: #262626;
}
</style>