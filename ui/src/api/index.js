import client from './client'

export default {
  login(loginForm) {
    return client.post('/login', loginForm);
  },
  signup(signupForm) {
    return client.post('/signup', signupForm);
  },
  getCaptcha() {
    return client.get('/captcha');
  },
  verifyCaptcha(captchaId, value) {
    return client.get(`/captcha/verify?captchaId=${captchaId}&value=${value}`);
  },
  getBooks() {
    return client.get('/books');
  },
  getParts(bookId) {
    return client.get(`/books/${bookId}/parts`);
  },
  getPart(partId) {
    return client.get(`/parts/${partId}`);
  },
  getPages(partId) {
    return client.get(`/parts/${partId}/pages`);
  },
  getUserInfo() {
    return client.get('/user');
  },
  searchUsers(search) {
    return client.post('/users/search', search);
  },
  updateUser(user) {
    return client.put(`/user/edit`, user);
  },
  viewPage(pageId, editable) {
    return client.get(`/pages/${pageId}?editable=${editable}`);
  },
  getHistoricalPages(pageId) {
    return client.get(`/pages/${pageId}/history`);
  },
  getTags() {
    return client.get(`/tags`);
  },
  addBook(book) {
    return client.post('/books', book);
  },
  addPart(part) {
    return client.post('/parts', part);
  },
  addPage(page) {
    return client.post('/pages', page);
  },
  editBook(book) {
    return client.put(`/books/${book.ID}`, book);
  },
  editPart(part) {
    return client.put(`/parts/${part.ID}`, part);
  },
  updatePage(page) {
    return client.put(`/pages/${page.ID}`, page);
  },
  editPage(page) {
    return client.put(`/pages/${page.ID}/edit`, page);
  },
  deleteBook(bookId) {
    return client.delete(`/books/${bookId}`);
  },
  deletePart(partId) {
    return client.delete(`/parts/${partId}`);
  },
  deletePage(pageId) {
    return client.delete(`/pages/${pageId}`);
  },
  shareBook(bookShare) {
    return client.post('/books/share', bookShare);
  },
  cancelShareBook(bookShare) {
    return client.post('/books/share/cancel', bookShare);
  },
  getBookSharedUsers(bookId) {
    return client.get(`/users/shared/book/${bookId}`);
  },
  modifyPassword(passwordModify) {
    return client.put('/user/modify-password', passwordModify);
  },
  siteSearch(keyword) {
    return client.get(`/site/search?keyword=${keyword}`);
  },
  sortBooks(sortedBooks) {
    return client.post(`/books/sort`, sortedBooks);
  },
  sortParts(sortedParts) {
    return client.post(`/parts/sort`, sortedParts);
  },
  sortPages(sortedPages) {
    return client.post(`/pages/sort`, sortedPages);
  },
  toggleStarBook(bookId) {
    return client.put(`/books/${bookId}/star`);
  },
  toggleStarPart(partId) {
    return client.put(`/parts/${partId}/star`);
  },
  toggleStarPage(pageId) {
    return client.put(`/pages/${pageId}/star`);
  },
  getSharedBooks() {
    return client.get('/shared/books');
  },
  getSharedParts(bookId) {
    return client.get(`/shared/books/${bookId}/parts`);
  },
  getSharedPages(bookId, partId) {
    return client.get(`/shared/books/${bookId}/parts/${partId}/pages`);
  },
  getSharedPage(bookId, partId, pageId) {
    return client.get(`/shared/books/${bookId}/parts/${partId}/pages/${pageId}`);
  },
  uploadFile(file) {
    let formdata = new FormData();
    formdata.append('file', file);
    return client.post('/upload', formdata, {'Content-Type': 'multipart/form-data'});
  }
}