import React from "react";
import styles from "./index.less";
import {Typography} from "antd";

const CodePreview: React.FC = ({children}) => (
  <pre className={styles.pre}>
    <code>
      <Typography.Text copyable>{children}</Typography.Text>
    </code>
  </pre>
);

export default CodePreview
