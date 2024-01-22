import React from "react";
import {PageContainer} from "@ant-design/pro-components";

const Setting: React.FC = () => {
  return <PageContainer>
    用户信息：包括后台账户密码修改。 <br />
    server-node: 已部署的node: <br />
    - tunnel-cloud 云厂商(ali、tencent、huawei、baidu、sealos)  <br />
    - tunnel-position 地域  <br />
    - tunnel-type 隧道类型  <br />
    - tunnel-status 隧道状态  <br />
    - tunnel-speed 隧道传输速率  <br />
    - tunnel-billing 账单花费 <br />
  </PageContainer>
}

export default Setting
