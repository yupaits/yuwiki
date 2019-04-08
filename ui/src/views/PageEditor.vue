<template>
  <div>
    <a-row class="opt-holder">
      <a-col :span="16">
        <span>
          <a-button icon="rollback" class="ml-1" @click="$router.go(-1)">返回</a-button>
          <span class="page-title ml-3">{{viewedPage.title}}</span>
        </span>
      </a-col>
      <a-col :span="4">
        <span>
          <a-input placeholder="请选择或填写标签"></a-input>
        </span>
      </a-col>
      <a-col :span="4">
        <span class="pull-right">
          <a-button icon="save" class="mr-1">存为草稿</a-button>
          <a-button type="primary" icon="check-circle" class="mr-1">发布页面</a-button>
        </span>
      </a-col>
    </a-row>
    <mavon-editor :value="viewedPage.content" :toolbars="toolbars" class="editor"></mavon-editor>
  </div>
</template>

<script>
import config from '../config'
export default {
  computed: {
    viewedPage() {
      return this.$store.getters.record;
    }
  },
  data() {
    return {
      toolbars: config.editor.toolbars
    }
  }
}
</script>

<style scoped>
.opt-holder {
  line-height: 60px;
}
.page-title {
  font-size: 20px;
  font-weight: bold;
}
.editor {
  z-index: 0;
  margin: 8px;
  height: calc(100vh - 78px);
}
</style>
