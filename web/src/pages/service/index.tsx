import React from "react";
import {PageContainer} from "@ant-design/pro-components";

const LocalService: React.FC = () => {
  return <PageContainer
    title={"本地服务"}
    content={"运行在本地的一些服务，用于处理一些云端数据的解析客户端"}
    extra={""} />
}

export default LocalService
