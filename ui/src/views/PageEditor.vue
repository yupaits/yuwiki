<template>
  <div class="editor-page">
    <a-row :gutter="16">
      <a-col :span="4">
        <div class="editor-header">
          <a-button icon="rollback" @click="$router.go(-1)">返回</a-button>
        </div>
        <div>
          <h4>页面标签：</h4>
          <a-input placeholder="请选择或填写标签"></a-input>
        </div>
        <div class="mt-3">
          <h4>页面历史：</h4>
          <a-spin :spinning="historyLoading">
            <div v-if="noHistory">
              此页面没有历史记录
            </div>
            <div v-else>
              <div v-for="page in historicalPages" :key="page.ID" class="history-item" :class="{'active': historyId === page.ID}" @click="showHistory(page)">
                <a-icon type="edit"/> 发布于 {{dayjs(page.CreatedAt).format('YYYY-MM-DD HH:mm:ss')}}
              </div>
            </div>
          </a-spin>
        </div>
      </a-col>
      <a-col :span="20">
        <a-row class="editor-header">
          <a-col :span="20">
            <h3 class="page-title">{{viewedPage.title}}</h3>
          </a-col>
          <a-col :span="4">
            <span class="pull-right">
              <a-button icon="save" class="mr-1" @click="saveDraft">存为草稿</a-button>
              <a-button type="primary" icon="check-circle" @click="publish">发布页面</a-button>
            </span>
          </a-col>
        </a-row>
        <mavon-editor :value="viewedPage.content" :toolbars="toolbars" ref="editor" class="editor"></mavon-editor>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import config from '../config'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')

export default {
  computed: {
    viewedPage() {
      return this.$store.getters.record;
    }
  },
  data() {
    return {
      dayjs,
      toolbars: config.editor.toolbars,
      historyLoading: false,
      historicalPages: [],
      noHistory: false,
      historyId: undefined
    }
  },
  created() {
    this.fetchHistoricalPages();
  },
  methods: {
    fetchHistoricalPages() {
      this.historyLoading = true;
      const pageId = this.$store.getters.pageId;
      if (pageId) {
        this.$api.getHistoricalPages(pageId).then(res => {
          const historicalPages = res.data;
          this.historicalPages = res.data;
          if (historicalPages.length === 0) {
            this.noHistory = true;
          } else {
            this.noHistory = false;
          }
          this.historyLoading = false;
        }).catch(() => {
          this.historyLoading = false;
        });
      } else {
        this.$router.push('/');
      }
    },
    showHistory(page) {
      if (this.historyId === page.ID) {
        this.$refs.editor.d_value = this.viewedPage.content;
        this.historyId = undefined;
      } else {
        this.$refs.editor.d_value = page.Content;
        this.historyId = page.ID;
      }
    },
    saveDraft() {
      this.savePage(false);
    },
    publish() {
      this.savePage(true);
    },
    savePage(published) {
      this.viewedPage.content = this.$refs.editor.d_value;
      this.viewedPage.published = published;
      this.$api.editPage(this.viewedPage).then(() => {
        this.$message.success(published ? '发布页面成功' : '存为草稿成功');
        this.$router.push('/');
      });
    }
  }
}
</script>

<style scoped>
.editor-page {
  padding: 0 16px;
}
.editor-header {
  line-height: 60px;
}
.page-title {
  font-size: 20px;
  font-weight: bold;
}
.editor {
  z-index: 0;
  height: calc(100vh - 78px);
}
.history-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.history-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.history-item.active {
  background: #91d5ff;
  font-weight: bold;
  color: #262626;
}
</style>
