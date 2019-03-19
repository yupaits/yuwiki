<template>
  <div>
    <a-form>
      <a-form-item label="笔记本" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-select v-model="part.bookId" placeholder="请选择笔记本"></a-select>
      </a-form-item>
      <a-form-item label="分区组" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <a-select v-model="part.groupId" placeholder="请选择所属分区组" allowClear></a-select>
      </a-form-item>
      <a-form-item label="分区名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input v-model="part.name" placeholder="请填写分区名"></a-input>
      </a-form-item>
      <a-form-item label="分区类型" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-radio-group v-model="partType" @change="handlePartTypeChange">
          <a-radio value="PART">分区</a-radio>
          <a-radio value="GROUP">分区组</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item label="密码保护" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-switch v-model="part.protected"></a-switch>
        <a-input type="password" v-model="part.password" placeholder="请填写保护密码" v-if="part.protected"></a-input>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
export default {
  computed: {
    part() {
      return this.$store.getters.record;
    }
  },
  watch: {
    part(val) {
      this.partType = val.partType || 'PART';
      this.$store.dispatch('setRecord', val);
    }
  },
  data() {
    return {
      partType: 'PART'
    }
  },
  methods: {
    handlePartTypeChange(e) {
      const partType = e.target.value;
      this.partType = partType;
      this.part.partType = partType;
    }
  }
}
</script>

<style scoped>

</style>
