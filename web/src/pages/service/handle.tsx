import {FormValueType} from "@/pages/service/components/CreateForm";
import {message} from "antd";
import {createServiceProxy, deleteServiceProxy, speedServiceProxy, updateServiceProxy} from "@/services/service/api";

export const handleCreateProxy = async (fields: FormValueType) => {
  const hide = message.loading('创建中......');
  try {
    // 处理一下参数
    fields.status = 1;
    const {success, code, msg} = await createServiceProxy(fields);
    hide();
    if (success) {
      message.success('创建成功');
    } else {
      message.error("创建失败: " +
        code + ":" + msg
      )
    }
  } catch (error) {
    hide();
    message.error('创建失败: ' + error);
  }
};

export const handleUpdateProxy = async (data: FormValueType) => {
  const hide = message.loading('更新中......');
  try {
    const {success, code, msg} = await updateServiceProxy(data);
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


export const handleDeleteProxy = async (data: FormValueType) => {
  const hide = message.loading('删除中......');
  try {
    const {success, code, msg} = await deleteServiceProxy(data);
    hide();
    if (success) {
      message.success('删除成功');
    } else {
      message.error("删除失败: " +
        code + ":" + msg
      )
    }
  } catch (error) {
    hide();
    message.error('删除失败:' + error);
  }
}

export const handleSpeedProxy = async (data: FormValueType) => {
  const hide = message.loading('测速中......');
  try {
    const {success, code, msg} = await speedServiceProxy(data);
    hide();
    if (!success) {
      message.error("测速失败: " +
        code + ":" + msg
      )
    }
  } catch (error) {
    hide();
    message.error('测速失败:' + error);
  }
}
