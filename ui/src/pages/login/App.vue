<template>
  <div class="login-page">
    <a-row type="flex" justify="center">
      <a-col>
        <a-card class="login-card">
          <div slot="title">
            <img src="/favicon.ico" alt="logo" height="32" width="32">
            <span class="ml-1">登录</span>
          </div>
          <a-form>
            <a-form-item hasFeedback :validateStatus="validate.username.status" :help="validate.username.help">
              <a-input v-model="loginForm.username" placeholder="请填写用户名" @keyup.enter="login">
                <span slot="prefix">
                  <a-icon type="user"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item hasFeedback :validateStatus="validate.password.status" :help="validate.password.help">
              <a-input type="password" v-model="loginForm.password" placeholder="请填写登录密码" @keyup.enter="login">
                <span slot="prefix">
                  <a-icon type="lock"/>
                </span>
              </a-input>
            </a-form-item>
            <a-form-item class="is-marginless">
              <a-button type="primary" block @click="login">登录</a-button>
              <a href="/signup">现在注册！</a>
              <a class="login-forget" href="">忘记密码</a>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      loginForm: {
        username: '',
        password: '',
      },
      validate: {
        username: {
          status: '',
          help: ''
        },
        password: {
          status: '',
          help: ''
        }
      }
    }
  },
  methods: {
    login() {
      if (this.validateForm()) {
        this.$api.login(this.loginForm).then(() => {
          window.location.replace('/');
        });
      }
    },
    validateForm() {
      let result = true;
      if (!this.loginForm.username) {
        this.validate.username = {status: 'error', help: '用户名不能为空！'};
        result = false;
      } else {
        this.validate.username = {status: 'success', help: ''};
      }
      if (!this.loginForm.password) {
        this.validate.password = {status: 'error', help: '登录密码不能为空！'};
        result = false;
      } else if (this.loginForm.password.length < 6) {
        this.validate.password = {status: 'error', help: '登录密码长度不能小于6！'};
        result = false;
      } else {
        this.validate.password = {status: 'success', help: ''};
      }
      return result;
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
}
.login-card {
  background: rgba(255, 255, 255, .90);
  margin-top: 5rem;
  width: 25rem;
}
.login-forget {
  float: right;
}
.other-tip {
  font-size: 14px;
}
</style>