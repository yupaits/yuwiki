<template>
  <div class="profile-page">
    <a-row>
      <a-col :span="8" :offset="8">
        <h3>
          <span><a-icon type="user"/> 个人设置</span>
          <span>
            <a-button type="dashed" size="small" class="pull-right" icon="rollback" @click="$router.go(-1)">返回</a-button>
          </span>
        </h3>
        <a-card class="mt-2">
          <a-form>
            <a-form-item label="头像" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
              <a-avatar :size="64" shape="square" :src="user.avatar" icon="user" :loadError="() => false"></a-avatar>
              <a-input v-model="user.avatar" placeholder="请输入头像链接" class="mt-1"></a-input>
            </a-form-item>
            <a-form-item label="用户名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
              <span>{{user.username}}</span>
            </a-form-item>
            <a-form-item label="姓名" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
              <a-input v-model="user.name" placeholder="请输入姓名"></a-input>
            </a-form-item>
            <a-form-item label="电子邮箱" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
              <a-input v-model="user.email" placeholder="请填写邮箱地址"></a-input>
            </a-form-item>
            <a-form-item label="手机号码" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper" required>
              <a-input v-model="user.phone" placeholder="请输入手机号码"></a-input>
            </a-form-item>
            <a-form-item label="生日" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
              <a-date-picker :value="birthday" :disabledDate="disabledDate" allowClear format="YYYY-MM-DD" @change="handleBirthdayChange"></a-date-picker>
            </a-form-item>
            <a-form-item label="性别" :labelCol="$styles.form.label" :wrapperCol="$styles.form.wrapper">
              <a-radio-group v-model="user.gender">
                <a-radio-button v-for="(gender, name) in $messages.enums.gender" :key="name" :value="name">
                  <a-icon :type="gender.icon"/> {{gender.label}}
                </a-radio-button>
              </a-radio-group>
            </a-form-item>
            <a-form-item :wrapperCol="$styles.form.contentWrapper">
              <a-button type="primary" @click="saveProfile">确认修改</a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import moment from 'moment'
export default {
  computed: {
    user() {
      return this.$store.getters.user;
    },
    birthday() {
      return this.user.birthday ? moment(this.user.birthday) : undefined;
    }
  },
  methods: {
    handleBirthdayChange(date, dateString) {
      this.$set(this.user, 'birthday', dateString);
    },
    disabledDate(current) {
      return current >= moment().endOf('day');
    },
    saveProfile() {
      this.$api.updateUser(this.user).then(() => {
        this.$store.dispatch('setUser', this.user);
        this.$message.success(this.$messages.result.updateSuccess);
      });
    }
  }
}
</script>

<style scoped>
.profile-page {
  padding: 16px 24px;
}
.profile-form {
  border: 1px solid #eee;
}
</style>