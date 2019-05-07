<template>
  <div>
    <a-layout>
      <a-layout-header>
        <span class="logo" @click="home">
          <img src="favicon.ico" alt="YuWIki">
          <span class="ml-1">知识库</span>
        </span>
        <span class="pull-right">
          <a-button icon="rollback" @click="$router.go(-1)">返回</a-button>
        </span>
      </a-layout-header>
      <a-layout-content>
        <a-row class="search-container">
          <a-col :span="8">
            <a-input-search size="large" v-model="keyword" class="mt-3" @search="searchPages"></a-input-search>
          </a-col>
        </a-row>
        <a-divider />
        <a-row class="search-container">
          <a-col :span="8">
            <a-spin :spinning="loading">
              <div v-if="pages.length > 0">
                <div class="mb-1">找到了{{pages.length}}个结果</div>
                <div class="list">
                  <div v-for="page in pages" :key="page.ID" class="page-item" :class="{'active': pageId === page.ID}" @click="selectPage(page)">
                    <div class="text-primary page-title">{{page.title}}</div>
                    <div class="page-addition">
                      <a-icon type="clock-circle"/> 创建于 
                      <span :title="dayjs(page.CreatedAt).format('YYYY年MM月DD日 HH:mm:ss')">{{dayjs().from(dayjs(page.CreatedAt))}}</span>
                      <span class="ml-1" v-if="page.tags && page.tags.length > 0"><a-icon type="tags"/> {{page.tags.join(', ')}}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else>
                <div>没有与此相关的结果！</div>
                <div class="list"></div>
              </div>
            </a-spin>
          </a-col>
          <a-col :span="1"></a-col>
          <a-col :span="15">
            <div v-if="viewedPage.ID">
              <mavon-editor :value="this.viewedPage.content" :boxShadow="false" :toolbars="toolbars" :editable="false" defaultOpen="preview" :subfield="false" class="page-preview"></mavon-editor>
            </div>
          </a-col>
        </a-row>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script>
import config from '../config'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')
dayjs.extend(relativeTime)

export default {
  data() {
    return {
      keyword: '',
      loading: false,
      pages: [],
      pageId: undefined,
      viewedPage: {},
      dayjs,
      toolbars: config.preivew.toolbars
    }
  },
  created() {
    const keyword = this.$store.getters.keyword;
    if (keyword) {
      this.keyword = keyword;
      this.searchPages();
    }
  },
  methods: {
    searchPages() {
      this.loading = true;
      this.viewedPage = {};
      this.pageId = undefined;
      this.$api.siteSearch(this.keyword).then(res => {
        this.pages = res.data || [];
        this.loading = false;
      }).catch(() => {
        this.loading = false;
      });
    },
    viewPage() {
      this.$api.viewPage(this.pageId, false).then(res => {
        this.viewedPage = res.data;
      });
    },
    home() {
      this.$router.push('/');
    },
    selectPage(page) {
      const pageId = page.ID;
      if (pageId) {
        this.pageId = pageId;
        this.viewPage();
      }
    }
  }
}
</script>

<style scoped>
::-webkit-scrollbar
{
  width: 6px;
  background-color: #f5f5f5;
}
::-webkit-scrollbar-track
{
  border-radius: 10px;
  background-color: #f5f5f5;
}
::-webkit-scrollbar-thumb
{
  border-radius: 10px;
  background-color: #bfbfbf;
}
.ant-layout {
  background: #fff;
}
.ant-layout-header {
  background: #f0f2f5;
}
.logo {
  cursor: pointer;
  font-size: 24px;
  line-height: 64px;
}
.search-container {
  padding: 0 50px;
}
.list {
  height: calc(100vh - 230px);
  overflow-x: hidden;
  overflow-y: auto;
}
.page-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.page-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.page-item.active {
  background: #91d5ff;
  font-weight: bold;
  color: #262626;
}
.page-title {
  font-size: 18px;
}
.page-preview {
  height: calc(100vh - 230px);
  z-index: 0;
}
</style>
