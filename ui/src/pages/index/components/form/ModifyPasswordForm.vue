<template>
  <div>
    <a-form>
      <a-form-item label="当前密码" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input type="password" v-model="passwordModify.oldPassword" placeholder="请输入当前密码"></a-input>
      </a-form-item>
      <a-form-item label="新的密码" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input type="password" v-model="passwordModify.newPassword" placeholder="请输入新的密码"></a-input>
      </a-form-item>
      <a-form-item label="确认密码" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
        <a-input type="password" v-model="passwordModify.confirmPassword" placeholder="请再次输入以确认密码"></a-input>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
export default {
  props: {
    passwordModify: {
      type: Object,
      requried: true
    }
  },
  watch: {
    passwordModify(modify) {
      this.$emit('update', modify);
    }
  }
}
</script>

<style scoped>

</style>
