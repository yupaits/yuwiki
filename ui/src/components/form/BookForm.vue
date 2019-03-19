<template>
  <div>
    <a-form>
      <a-form-item label="笔记本名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <a-input v-model="record.name" placeholder="请填写笔记本名"></a-input>
      </a-form-item>
      <a-form-item label="颜色" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
        <a-input v-model="record.name"></a-input>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
export default {
  props: {
    record: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      book: this.record
    }
  }
}
</script>

<style scoped>

</style>
