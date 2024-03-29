import {BookOutlined, GithubOutlined} from '@ant-design/icons';
import { Space } from 'antd';
import React from 'react';
import { useModel } from 'umi';
import styles from './index.less';

export type SiderTheme = 'light' | 'dark';

const GlobalHeaderRight: React.FC = () => {
  const { initialState } = useModel('@@initialState');

  if (!initialState || !initialState.settings) {
    return null;
  }

  const { navTheme, layout } = initialState.settings;
  let className = styles.right;

  if ((navTheme === 'dark' && layout === 'top') || layout === 'mix') {
    className = `${styles.right}  ${styles.dark}`;
  }
  return (
    <Space className={className}>
      <span
        className={styles.action}
        onClick={() => {
          window.open('https://seamoon.dvkunion.cn');
        }}
      >
        <BookOutlined />
      </span>
      <span
        className={styles.action}
        onClick={() => {
          window.open('https://www.github.com/DVKunion/Seamoon');
        }}>
        <GithubOutlined />
      </span>
      {/*<Avatar />*/}
    </Space>
  );
};
export default GlobalHeaderRight;
