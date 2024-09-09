import React, {useState} from "react";
import {PageContainer} from "@ant-design/pro-components";
import ProviderTable from "@/pages/account/provider/Provider";
import AdminTable from "@/pages/account/admin/Admin";
import TunnelTable from "@/pages/account/tunnel/Tunnel";

const Account: React.FC = () => {

  const [currentTab, setCurrentTab] = useState<string>("cloud_provider");

  return <PageContainer
    title={"账户管理"}
    tabList={[
      {
        tab: '云账户',
        key: 'cloud_provider',
      },
      {
        tab: '隧道账户',
        key: 'proxy_auth',
      },
      {
        tab: '系统账户',
        key: 'admin_auth',
      }
    ]}
    onTabChange={
      (key) => {
        setCurrentTab(key)
      }
    }
    extra={""}>
    {
      currentTab == "admin_auth" ? <AdminTable/> :
        currentTab == "proxy_auth" ? <TunnelTable/> : <ProviderTable/>
    }
  </PageContainer>
}

export default Account
