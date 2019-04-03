<template>
  <div>
    <div class="part-item" :class="{'active': $store.getters.partId === option.ID}" @click="selectPart(option.ID)">
      <a-icon :type="option.partType === 0 ? 'folder-open' : 'inbox'"/> {{option.name}}
    </div>
    <part-node v-for="part in option.SubParts" :key="part.ID" :option="part" class="ml-2"/>
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
    },
    methods: {
      selectPart(partId) {
        this.$store.dispatch('setPartId', partId);
        let part = Object.assign({}, this.option);
        delete part.SubParts;
        this.$store.dispatch('setPart', part);
      }
    }
  }
</script>

<style scoped>
.part-item {
  font-size: 16px;
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