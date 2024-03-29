import {message} from "antd";
import {updateSysConfig} from "@/services/setting/api";

export const handleUpdateSysConfig = async (data: Config.SystemConfig) => {
  const hide = message.loading('更新中......');
  try {
    if (data.auto_start) {
      data.auto_start = "true"
    } else {
      data.auto_start = "false"
    }
    const {success, code, msg} = await updateSysConfig(data);
    hide();
    if (success) {
      message.success('更新成功');
    } else {
      message.error("更新失败: " +
        code + ":" + msg
      )
    }
  } catch (error) {
    hide();
    message.error('更新失败:' + error);
  }
}
