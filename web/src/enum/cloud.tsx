import {Space} from "antd";
import IconFont from "@/components/IconFont";

export const CloudProvideTypeIcon = [
  <></>,
  <IconFont type={"icon-aliyun"}/>,
  <IconFont type={"icon-tengxunyun1"}/>,
  <IconFont type={"icon-huaweiyun1"}/>,
  <IconFont type={"icon-baiduyun"}/>,
  <IconFont type={"icon-sealos"}/>
]

export const CloudProvideTypeValueEnum = {
  1: <Space><IconFont type={"icon-aliyun"}/>阿里云</Space>,
  2: <Space><IconFont type={"icon-tengxunyun1"}/>腾讯云</Space>,
  // 3: <Space><IconFont type={"icon-huaweiyun1"}/>华为云</Space>,
  // 4: <Space><IconFont type={"icon-baiduyun"}/>百度云</Space>,
  5: <Space><IconFont type={"icon-sealos"}/>Sealos</Space>,
}


export const CloudProvideTypeEnum = {
  1: {
    text: <Space><IconFont type={"icon-aliyun"}/>阿里云</Space>,
  },
  2: {
    text: <Space><IconFont type={"icon-tengxunyun1"}/>腾讯云</Space>,
  },
  // 3: {
  //   text: <Space><IconFont type={"icon-huaweiyun1"}/>华为云</Space>,
  // },
  // 4: {
  //   text: <Space><IconFont type={"icon-baiduyun"}/>百度云</Space>,
  // },
  5: {
    text: <Space><IconFont type={"icon-sealos"}/>Sealos</Space>,
  },
}

export const CloudProviderStatusEnum = {
  1: {
    text: '创建中',
    status: 'processing',
  },
  2: {
    text: '正常',
    status: 'success',
  },
  3: {
    text: '异常',
    status: 'error',
  },
  4: {
    text: '同步中',
    status: 'default',
  },
  5: {
    text: '已禁用',
    status: 'error',
  },
  6: {
    text: '同步失败',
    status: 'error',
  },
  7: {
    text: '删除中',
    status: 'warning',
  }
}

export const ALiYunRegionEnum = {
  // 阿里云
  "cn-hangzhou": "华东1(杭州)",
  "cn-shanghai": "华东2(上海)",
  "cn-qingdao": "华北1(青岛)",
  "cn-beijing": "华北2(北京)",
  "cn-zhangjiakou": "华北3(张家口)",
  "cn-huhehaote": "华北5(呼和浩特)",
  "cn-shenzhen": "华南1(深圳)",
  "cn-chengdu": "西南1(成都)",
  "cn-hongkong": "中国香港",
  "ap-northeast-1": "日本(东京)",
  // "ap-northeast-2": "韩国(首尔)",
  "ap-southeast-1": "新加坡(新加坡)",
  "ap-southeast-2": "澳大利亚(悉尼)",
  "ap-southeast-3": "马来西亚(吉隆坡)",
  "ap-southeast-5": "印尼(雅加达)",
  // "ap-southeast-7": "泰国(曼谷)",
  "ap-south-1": "印度(孟买)",
  "eu-central-1": "德国(法兰克福)",
  "eu-west-1": "英国(伦敦)",
  "us-west-1": "美国(硅谷)",
  "us-east-1": "美国(弗吉尼亚)",
}

export const SealosRegionEnum = {
  "beijing-a":   "北京 A",
  "singapore-b": "新加坡 B",
  "guangzhou-g": "广州 G",
  "hangzhou-h":  "杭州 H",
}

export const TencentRegionEnum = {
  "ap-beijing": "华北(北京)",
  "ap-chengdu": "西南(成都)",
  "ap-guangzhou": "华南(广州)",
  "ap-shanghai": "华东(上海)",
  "ap-nanjing": "华东(南京)",
  "ap-hongkong": "中国香港",
  "ap-mumbai": "亚太南部(孟买)",
  "ap-singapore": "亚太东南(新加坡)",
  "ap-bangkok": "亚太东南(曼谷)",
  "ap-seoul": "亚太东北(首尔)",
  "ap-tokyo": "亚太东北(东京)",
  "eu-frankfurt": "欧洲(法兰克福)",
  "na-ashburn": "美国东部(弗吉尼亚)",
  // "na-toronto": "北美(多伦多)",
  "na-siliconvalley": "美国西部(硅谷)",
}

export const RegionEnum = {
  ...ALiYunRegionEnum,
  ...TencentRegionEnum,
  ...SealosRegionEnum,
}
