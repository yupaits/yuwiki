import client from './client'

export default {
  getBooks() {
    return client.get('/books');
  },
  getParts(bookId) {
    return client.get(`/books/${bookId}/parts`);
  },
  getPages(partId) {
    return client.get(`/parts/${partId}/pages`);
  },
  getStarItems() {
    return client.get('/books/star');
  },
  getSharedBooks() {
    return client.get('/books/shared');
  },
  getUserInfo() {
    return client.get('/user');
  },
  viewPage(pageId) {
    return client.get(`/pages/${pageId}`);
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
    return client.put(`/books/${book.id}`, book);
  },
  editPart(part) {
    return client.put(`/parts/${part.id}`, part);
  },
  editPage(page) {
    return client.put(`/pages/${page.id}`, page);
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
    return client.put('/user/modify-passwd', passwordModify);
  },
  siteSearch(query) {
    return client.post('/site/search', query);
  }
}