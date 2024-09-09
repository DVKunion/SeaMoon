import React from "react";
import {PageContainer} from "@ant-design/pro-components";

const Application: React.FC = () => {
  return <PageContainer
    title={"应用函数"}
    tabList={[
      {
        tab: '临时邮箱',
        key: 'email',
      },
      {
        tab: '反弹shell',
        key: 'shell',
      },
    ]}
    extra={""} />
}

export default Application
