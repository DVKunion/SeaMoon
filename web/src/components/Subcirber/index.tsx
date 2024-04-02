import {MenuProps, message, Button, Dropdown} from "antd";
import React from "react";
import IconFont from "@/components/IconFont";

export const SSrDropDown: React.FC = () => {
  const items: MenuProps['items'] = [
    {
      label: (<a
        href={"sub://" + btoa(window.location.host + "/api/v1/tunnel/subscribe/ss/?token=" + localStorage.getItem("token"))}>一键导入
        SS 客户端</a>),
      key: '0',
    },
    {
      label: <a onClick={() => {
        const dir = "http://" + window.location.host + "/api/v1/tunnel/subscribe/ss/?token=" + localStorage.getItem("token");
        navigator.clipboard.writeText(dir).then(() => {
          message.success("复制 SS 订阅地址成功")
        })
      }}>手动复制 SS 订阅地址</a>,
      key: '1',
    },
  ];
  return <Dropdown menu={{items}} arrow>
    <Button shape={"round"} ghost icon={<IconFont type={"icon-Shadowsocks-Logo-copy"}/>}>SS</Button>
  </Dropdown>
}

export const ClashDropDown: React.FC = () => {
  const items: MenuProps['items'] = [
    {
      label: (<a
        href={"clash://install-config?url=http://" + window.location.host + "/api/v1/tunnel/subscribe/clash/?token=" + localStorage.getItem("token")}>一键导入
        Clash 客户端</a>),
      key: '0',
    },
    {
      label: <a onClick={() => {
        const dir = "http://" + window.location.host + "/api/v1/tunnel/subscribe/clash/?token=" + localStorage.getItem("token");
        navigator.clipboard.writeText(dir).then(() => {
          message.success("复制 Clash 订阅地址成功")
        })
      }}>手动复制 Clash 订阅地址</a>,
      key: '1',
    },
  ];
  return <Dropdown menu={{items}} arrow>
    <Button shape={"round"} ghost
            icon={<IconFont type={"icon-weibiaoti-_huabanfuben-copy-copy-copy-copy"}/>}>Clash</Button>
  </Dropdown>
}


export const ShadowRocketDropDown: React.FC = () => {
  const items: MenuProps['items'] = [
    {
      label: (<a
        href={"shadowrocket://add/sub://" + Buffer.from("http://"+ window.location.host + "/api/v1/tunnel/subscribe/shadowrocket/?token=" + localStorage.getItem("token"), 'utf8').toString('base64')}>一键导入
        ShadowRocket 客户端</a>),
      key: '0',
    },
    {
      label: <a onClick={() => {
        const dir = "http://" + window.location.host + "/api/v1/tunnel/subscribe/shadowrocket/?token=" + localStorage.getItem("token");
        navigator.clipboard.writeText(dir).then(() => {
          message.success("复制 ShadowRocket 订阅地址成功")
        })
      }}>手动复制 ShadowRocket 订阅地址</a>,
      key: '1',
    },
  ];
  return <Dropdown menu={{items}} arrow>
    <Button shape={"round"} ghost icon={<IconFont type={"icon-Untitled-"}/>}>ShadowRocket</Button>
  </Dropdown>
}

