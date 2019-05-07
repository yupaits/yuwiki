<template>
  <div class="share-page">
    <a-row>
      <a-col :span="8" :offset="8">
        <h3>
          <span><a-icon type="share-alt"/> 共享笔记本</span>
          <span>
            <a-button type="dashed" size="small" class="pull-right" icon="rollback" @click="$router.go(-1)">返回</a-button>
          </span>
        </h3>
        <a-card class="books-card mt-2">
          <a-spin :spinning="loading">
            <div v-if="books.length > 0">
              <div v-for="book in books" :key="book.ID" class="book-item">
                <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
                <span class="pull-right">
                  <a-button size="small" icon="share-alt" @click="shareBook(book)">共享</a-button>
                </span>
              </div>
            </div>
            <div v-else>
              <a-alert message="当前无笔记本，请先创建笔记本再进行共享！"></a-alert>
            </div>
          </a-spin>
        </a-card>
      </a-col>
    </a-row>

    <a-modal :visible.sync="shareVisible" :maskClosable="false" :footer="null" @cancel="shareVisible = false">
      <div slot="title">
        分享笔记本 <a-icon type="book" theme="twoTone" :twoToneColor="book.color"/> {{book.name}}
      </div>
      <div class="user-holder">
        <a-form layout="inline">
          <a-form-item label="选择目标用户">
            <a-select v-model="userId" placeholder="请输入用户名" showSearch allowClear :defaultActiveFirstOption="false" :filterOption="false" 
                @search="searchUser" @change="handleUserChange" class="user-select">
              <a-spin v-if="searching" slot="notFoundContent" size="small"></a-spin>
              <a-select-option v-for="user in users" :key="user.ID">
                <a-avatar size="small" shape="square" icon="user" :src="user.avatar"></a-avatar> {{user.nickname || user.username}}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" icon="check" @click="handleShareBook">分享</a-button>
          </a-form-item>
        </a-form>
        <a-divider/>
        <a-form layout="inline">
          <a-form-item label="筛选共享用户">
            <a-input v-model="filterKeyword" class="user-filter" placeholder="请输入关键字">
              <a-icon slot="prefix" type="filter"/>
            </a-input>
          </a-form-item>
        </a-form>
        <div class="mt-3">
          <a-row :gutter="12">
            <a-col v-for="user in filteredSharedUsers" :key="user.userId" :span="8">
              <div class="user-item mb-1">
                <a-avatar size="large" shape="square" icon="user" :src="user.avatar" class="mr-1"></a-avatar> {{user.nickname || user.username}}
                <a-popconfirm title="确定取消共享笔记本给此用户吗？" placement="top" @confirm="handleUnshareBook(user.userId)">
                  <a-button size="small" shape="circle" icon="close" title="取消共享" class="pull-right"></a-button>
                </a-popconfirm>
              </div>
            </a-col>
          </a-row>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script>
import debounce from 'lodash/debounce'
export default {
  data() {
    return {
      books: [{}],
      loading: false,
      shareVisible: false,
      book: {},
      users: [],
      userId: undefined,
      userKeyword: undefined,
      filterKeyword: '',
      searching: false,
      sharedUsers: []
    }
  },
  created() {
    this.fetchBooks();
  },
  computed: {
    filteredSharedUsers() {
      if (this.sharedUsers && this.sharedUsers.length > 0) {
        return this.sharedUsers.filter((user) => {
          return user.username.indexOf(this.filterKeyword) !== -1 || user.nickname.indexOf(this.filterKeyword) !== -1;
        });
      } else {
        return [];
      }
    }
  },
  methods: {
    fetchBooks() {
      this.loading = true;
      this.$api.getBooks().then(res => {
        this.books = res.data;
        this.loading = false;
      }).catch(() => {
        this.loading = false;
      });
    },
    fetchSharedUsers(bookId) {
      this.$api.getBookSharedUsers(bookId).then(res => {
        this.sharedUsers = res.data;
      });
    },
    fetchUsers: debounce((self, keyword) => {
      self.$api.searchUsers({keyword, noSelf: true}).then(res => {
        self.users = res.data;
        self.searching = false;
      }).catch(() => {
        self.searching = false;
      });
    }, 800),
    shareBook(book) {
      this.book = JSON.parse(JSON.stringify(book));
      this.fetchSharedUsers(book.ID);
      this.shareVisible = true;
    },
    searchUser(keyword) {
      this.searching = true;
      this.fetchUsers(this, keyword);
    },
    handleUserChange(keyword) {
      this.userKeyword = keyword;
    },
    handleShareBook() {
      this.$api.shareBook({bookId: this.book.ID, userId: this.userId}).then(() => {
        this.userKeyword = undefined;
        this.fetchSharedUsers(this.book.ID);
      });
    },
    handleUnshareBook(userId) {
      const bookId = this.book.ID;
      this.$api.cancelShareBook({bookId, userId}).then(() => {
        this.fetchSharedUsers(bookId);
      });
    }
  }
}
</script>

<style scoped>
.share-page {
  padding: 16px 24px;
}
.books-card {
  min-height: 500px;
}
.book-item {
  line-height: 32px;
  font-size: 16px;
  padding: 0 8px;
  border-radius: 4px;
  margin-bottom: 2px;
}
.book-item:hover {
  background: #e6f7ff;
}
.user-select {
  width: 260px;
}
.user-filter {
  width: 260px;
}
.user-holder {
  min-height: 360px;
}
.user-item {
  border: 1px solid #e8e8e8;
  border-radius: 2px;
  padding: 3px;
}
.unshare-btn {
  cursor: pointer;
}
</style>
