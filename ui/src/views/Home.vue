<template>
  <div>
    <a-layout id="components-layout-demo-top" class="layout">
      <a-layout-header>
        <span class="logo">
          <img src="favicon.ico" alt="YuWIki">
          <span class="ml-1">知识库</span>
        </span>
        <span class="pull-right">
          <a-dropdown>
            <a-menu slot="overlay" @click="handleCreate">
              <a-menu-item key="book"><a-icon type="book"/>笔记本</a-menu-item>
              <a-menu-item key="part"><a-icon type="folder"/>分区</a-menu-item>
              <a-menu-item key="page"><a-icon type="file-text"/>页面</a-menu-item>
            </a-menu>
            <a-button size="large" icon="plus" class="ml-1">
              新建 <a-icon type="down"/>
            </a-button>
          </a-dropdown>
          <a-dropdown>
            <a-menu slot="overlay" @click="handleUserOpt">
              <a-menu-item key="share-book"><a-icon type="share-alt"/>共享笔记本</a-menu-item>
              <a-menu-item key="modify-passwd"><a-icon type="key"/>修改密码</a-menu-item>
              <a-menu-item key="logout"><a-icon type="logout"/>注销登录</a-menu-item>
            </a-menu>
            <a-button size="large" class="ml-1">
              <a-avatar size="small" shape="square" icon="user" src="https://avatars1.githubusercontent.com/u/12194490?s=40&v=4"></a-avatar>
              <span class="ml-1">用户</span> <a-icon type="down"/>
            </a-button>
          </a-dropdown>
        </span>
      </a-layout-header>
      <a-layout-content class="layout-content">
        <a-row :gutter="16">
          <a-col :span="3">
            <div class="holder">
              <h3 class="text-primary text-bold holder-header">
                <a-icon type="book"/> 笔记本
                <span class="pull-right" v-if="bookId">
                  <a-button size="small" icon="edit" class="mr-1" @click="editBook"></a-button>
                  <a-popconfirm title="确定删除此笔记本吗？" placement="right" @confirm="handleDeleteBook">
                    <a-button size="small" icon="delete"></a-button>
                  </a-popconfirm>
                </span>
              </h3>
              <div class="list">
                <div v-for="i in 10" :key="i" class="book-item" :class="{'active': bookId === i}" @click="selectBook(i)">
                  <a-icon type="book"/> 第{{i}}本书
                </div>
              </div>
            </div>
          </a-col>
          <a-col :span="4">
            <div class="holder">
              <h3 class="text-primary text-bold"><a-icon type="folder"/> 分区</h3>
            </div>
          </a-col>
          <a-col :span="6">
            <div class="holder">
              <h3 class="text-primary text-bold"><a-icon type="file-text"/> 页面</h3>
            </div>
          </a-col>
          <a-col :span="11">
            <div class="holder">
            </div>
          </a-col>
        </a-row>
      </a-layout-content>
      <a-layout-footer style="text-align: center">
        <b>YuWiki</b> ©2019 <b><a href="https://github.com/YupaiTS" target="_blank">YupaiTS</a></b> 版权所有
      </a-layout-footer>
    </a-layout>

    <a-modal :visible.sync="modal.visible" 
             @cancel="modal.visible = false">
      <template slot="title">
        <a-icon :type="modalType[modal.type] ? modalType[modal.type].icon : 'question'"/> {{modalTitle}}
      </template>
    </a-modal>
  </div>
</template>

<script>
export default {
  data() {
    return {
      books: [],
      parts: [],
      pages: [],
      bookId: undefined,
      partId: undefined,
      pageId: undefined,
      book: {},
      part: {},
      page: {},
      modalVisible: false,
      modal: {
        type: undefined,
        key: undefined,
        visible: false
      },
      modalType: {
        create: {label: '新建', icon: 'plus'},
        modify: {label: '修改', icon: 'edit'}
      },
      title: {
        book: '笔记本',
        part: '分区',
        page: '页面',
      },
    }
  },
  computed: {
    modalTitle() {
      const type = this.modalType[this.modal.type];
      return (type ? type.label : '') + this.title[this.modal.key];
    },
    selectedBook() {
      return this.books.filter(book => book.id === this.bookId)[0] || {};
    }
  },
  methods: {
    fetchBooks() {

    },
    showModal(type, key) {
      this.modal = {
        type: type,
        key: key,
        visible: true
      };
    },
    handleCreate({key}) {
      this.showModal('create', key);
    },
    selectBook(i) {
      this.bookId = i;
      this.partId = undefined;
      this.pageId = undefined;
    },
    editBook() {
      this.book = JSON.parse(JSON.stringify(this.selectedBook));
      this.modal = {
        type: 'modify',
        key: 'book',
        visible: true
      };
    },
    handleDeleteBook() {
      alert('删除' + this.bookId);
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
.logo {
  color: #f5f5f5;
  font-size: 24px;
  line-height: 64px;
}
.layout-content {
  height: calc(100vh - 133px);
  padding: 24px 50px;
}
.holder {
  height: calc(100vh - 158px);
  padding: 8px 16px;
  border: 1px solid #f5f5f5;
  border-radius: 4px;
  background: #fff;
}
.holder-header {
  line-height: 28px;
}
.list {
  height: calc(100vh - 210px);
  overflow-x: hidden;
  overflow-y: auto;
}
.book-item,.part-item,.page-item {
  line-height: 32px;
  padding: 0 8px;
  border-radius: 4px;
}
.book-item:hover,.part-item:hover,.page-item:hover {
  cursor: pointer;
  background: #e6f7ff;
}
.book-item.active,.part-item.active,.page-item.active {
  background: #91d5ff;
}
</style>
