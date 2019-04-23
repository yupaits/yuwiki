<template>
  <div class="singup-page">
    <a-row type="flex" justify="center">
      <a-col>
        <a-card class="singup-card">
          <div slot="title">
            <img src="/favicon.ico" alt="logo" height="32" width="32">
            <span class="ml-1">注册账号</span>
          </div>
          <a-form>
            <a-form-item hasFeedback :validateStatus="validate.username.status" :help="validate.username.help">
              <a-input v-model="signupForm.username" placeholder="请填写用户名" @keyup.enter="singup">
                <span slot="prefix">
                  <a-icon type="user"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item hasFeedback :validateStatus="validate.email.status" :help="validate.email.help">
              <a-input type="email" v-model="signupForm.email" placeholder="请填写邮箱地址" @keyup.enter="singup">
                <span slot="prefix">
                  <a-icon type="mail"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item hasFeedback :validateStatus="validate.password.status" :help="validate.password.help">
              <a-input type="password" v-model="signupForm.password" placeholder="请填写登录密码" @keyup.enter="singup">
                <span slot="prefix">
                  <a-icon type="lock"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item hasFeedback :validateStatus="validate.confirmPassword.status" :help="validate.confirmPassword.help">
              <a-input type="password" v-model="signupForm.confirmPassword" placeholder="确认登录密码" @keyup.enter="singup">
                <span slot="prefix">
                  <a-icon type="lock"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item class="is-marginless">
              <a-button type="primary" block @click="singup">注册</a-button>
              <a href="/login">前往登录！</a>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script>
const emailReg = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/;
export default {
  name: "App",
  data() {
    return {
      signupForm: {
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
      },
      validate: {
        username: {
          status: '',
          help: ''
        },
        email: {
          status: '',
          help: ''
        },
        password: {
          status: '',
          help: ''
        },
        confirmPassword: {
          status: '',
          help: ''
        }
      }
    }
  },
  methods: {
    singup() {
      if (this.validateForm()) {
        this.$api.signup(this.signupForm).then(() => {
          window.location.replace('/');
        });
      }
    },
    validateForm() {
      let result = true;
      if (!this.signupForm.username) {
        this.validate.username = {status: 'error', help: '用户名不能为空！'};
        result = false;
      } else {
        this.validate.username = {status: 'success', help: ''};
      }
      if (!this.signupForm.email) {
        this.validate.email = {status: 'error', help: '邮箱地址不能为空！'};
        result = false;
      } else if (!emailReg.test(this.signupForm.email)) {
        this.validate.email = {status: 'error', help: '邮箱地址格式不正确！'};
        result = false;
      } else {
        this.validate.email = {status: 'success', help: ''}
      }
      if (!this.signupForm.password) {
        this.validate.password = {status: 'error', help: '登录密码不能为空！'};
        result = false;
      } else if (this.signupForm.password.length < 6) {
        this.validate.password = {status: 'error', help: '登录密码长度不能小于6！'};
        result = false;
      } else {
        this.validate.password = {status: 'success', help: ''};
      }
      if (!this.signupForm.confirmPassword) {
        this.validate.confirmPassword = {status: 'error', help: '登录密码不能为空！'};
        result = false;
      } else if (this.signupForm.confirmPassword.length < 6) {
        this.validate.confirmPassword = {status: 'error', help: '登录密码长度不能小于6！'};
        result = false;
      } else if (this.signupForm.password !== this.signupForm.confirmPassword) {
        this.validate.confirmPassword = {status: 'error', help: '两次输入的密码不一致！'};
      } else {
        this.validate.confirmPassword = {status: 'success', help: ''};
      }
      return result;
    }
  }
}
</script>

<style scoped>
.singup-page {
  min-height: 100vh;
}
.singup-card {
  background: rgba(255, 255, 255, .90);
  margin-top: 5rem;
  width: 25rem;
}
.other-tip {
  font-size: 14px;
}
</style>