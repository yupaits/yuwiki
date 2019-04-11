<template>
  <div>
    <div class="part-item" :class="{'active': $store.getters.partId === option.ID}" @click="selectPart(option.ID)">
      <a-icon :type="option.partType === 0 ? 'folder-open' : 'inbox'"/> {{option.name}}
    </div>
    <draggable v-model="option.SubParts">
      <transition-group>
        <part-node v-for="part in option.SubParts" :key="part.ID" :option="part" class="ml-2"/>
      </transition-group>
    </draggable>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  name: "PartNode",
  components: {draggable},
  props: {
    option: {
      type: Object,
      required: true
    }
  },
  methods: {
    selectPart(partId) {
      this.$store.dispatch('setPartId', partId);
      let part = Object.assign({}, this.option);
      delete part.SubParts;
      this.$store.dispatch('setPart', part);
      this.$emit('select', partId);
      this.$eventBus.$emit('selectPart', partId);
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