import client from './client'

export default {
  login(loginForm) {
    return client.post('/login', loginForm);
  },
  signup(signupForm) {
    return client.post('/signup', signupForm);
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
  getStarItems() {
    return client.get('/star/items');
  },
  getSharedBooks() {
    return client.get('/books/shared');
  },
  getUserInfo() {
    return client.get('/user');
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
    return client.post(`/books/share`, bookShare);
  },
  modifyPassword(passwordModify) {
    return client.put('/user/modify-password', passwordModify);
  },
  siteSearch(query) {
    return client.post('/site/search', query);
  },
  sortBooks(sortedBooks) {
    return client.post(`/books/sort`, sortedBooks);
  },
  sortParts(sortedParts) {
    return client.post(`/parts/sort`, sortedParts);
  },
  sortPages(sortedPages) {
    return client.post(`/pages/sort`, sortedPages);
  }
}