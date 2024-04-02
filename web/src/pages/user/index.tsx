import {login} from '@/services/user/api';
import {
  LockOutlined,
  UserOutlined,
} from '@ant-design/icons';
import {
  LoginForm,
  ProFormText,
} from '@ant-design/pro-components';
import {message} from 'antd';
import React from 'react';
import {history, useModel} from 'umi';
import styles from './index.less';
import ShieldList from "@/components/ShieldList";

const Login: React.FC = () => {
  const {initialState, setInitialState} = useModel('@@initialState');

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s: any) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };

  const handleSubmit = async (values: Auth.User) => {
    try {
      // 登录
      const msg = await login({...values});
      if (msg.success) {
        message.success('登录成功！');
        localStorage.setItem("user", values.Username || "")
        localStorage.setItem("token", msg.data || "")
        await fetchUserInfo();
        /** 此方法会跳转到 redirect 参数所在的位置 */
        if (!history) return;
        const {query} = history.location;
        const {redirect} = query as { redirect: string };
        history.push(redirect || '/');
        return;
      } else {
        message.error(msg.code + ":" + msg.msg)
      }
    } catch (error) {
      message.error('登录失败，请重试！');
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <img className={styles.logoImg} src="/icon_black.svg" alt="logo" width="768"/>
        <LoginForm
          // @ts-ignore
          title={<ShieldList/>}
          subTitle={"🌕 月出于云却隐于海"}

          onFinish={async (values) => {
            await handleSubmit(values as Auth.User);
          }}
        >
          <>
            <ProFormText
              name="Username"
              fieldProps={{
                size: 'large',
                prefix: <UserOutlined className={styles.prefixIcon}/>,
              }}
              placeholder={'用户名'}
              rules={[
                {
                  required: true,
                  message: "请输入用户名!"
                },
              ]}
            />
            <ProFormText.Password
              name="Password"
              fieldProps={{
                size: 'large',
                prefix: <LockOutlined className={styles.prefixIcon}/>,
              }}
              placeholder='密码'
              rules={[
                {
                  required: true,
                  message: "请输入密码！"
                },
              ]}
            />
          </>
        </LoginForm>
      </div>
    </div>
  );
};

export default Login;
