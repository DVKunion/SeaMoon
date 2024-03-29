import {extend} from 'umi-request';
import {history} from 'umi';
import {message} from "antd";

const request = extend({
  errorHandler: function (error) {
    const {response} = error;

    if (response) {
      if (response.status === 401 || response.status === 403) {
        response.clone().json().then(({code, msg}) => {
          // 使用message.error展示错误消息
          message.error(code + ":" + msg || '未登陆', 1).then(r => {
            // 清理当前认证信息
            localStorage.clear();
            // 当状态码为401或403时，重定向到登录页面
            history.push('/user/login');
          });
        });
      } else {
        response.clone().json().then(({code, msg}) => {
          message.error("请求失败 - " + code + ": " + msg).then()
        });
      }
    } else {
      message.error("请求失败 - 未知错误： " + error).then()
    }
    // 向外抛出错误
    throw error;
  },
});


export default request
