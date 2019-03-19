import axios from 'axios'
import {message} from 'ant-design-vue'

const client = axios.create({
  baseURL: '',
  headers: {'X-Requested-With': 'XMLHttpRequest'}
})

client.interceptors.response.use(res => {
  if (res.status !== 200) {
    return Promise.reject(res);
  }
  if (res.data.success) {
    return Promise.resolve(res.data);
  }
  message.error(res.data.message);
  return Promise.reject(res);
}, error => {
  if (error.response === undefined || error.response === null) {
    message.error(error);
    return Promise.reject(error);
  }
  if (error.response.status === 401) {
    utils.http.goPage('/login');
  }
  if (error.response.status === 403) {
    message.error('请求失败，无访问权限');
    return Promise.reject(error);
  }
  if (error.response.data.code) {
    message.error(`${error.response.data.code} - ${error.response.data.message}`);
    return Promise.reject(error.response.data);
  }
  message.error(`${error} - ${error.response.statusText}`);
  return Promise.reject(error.response);
});

export default client
