import {message} from "antd";
import {toNumber} from "lodash";
import type {FormValueType} from "@/pages/provider/components/CreateForm";
import {createProvider, deleteProvider, syncProvider, updateProvider} from "@/services/cloud/api";

export const handleCreateCloud = async (fields: FormValueType) => {
  const hide = message.loading('创建中......', 30);
  try {
    // 处理一下参数
    fields.status = 1;
    fields.type = toNumber(fields.type);
    const {success} = await createProvider(fields);
    hide();
    if (success) {
      message.success('创建成功');
      return true;
    }
  } catch (error) {
    hide();
  }
  return false;
};

export const handleUpdateCloud = async (fields: FormValueType) => {
  const hide = message.loading('更新中......');
  try {
    // 更新后，状态重至为0, 需要手动去同步认证
    fields.status = 1;
    const {success} = await updateProvider(fields);
    hide();
    if (success) {
      message.success('更新成功');
      return true;
    }
  } catch (error) {
    hide();
  }
  return false;
};

export const handleDeleteCloud = async (fields: number | undefined) => {
  const hide = message.loading('删除中......');
  if (fields === undefined) {
    hide();
    message.error("数据错误:不存在的数据记录")
    return false;
  }
  try {
    const {success} = await deleteProvider(fields);
    hide();
    if (success) {
      message.success('删除成功');
      return true;
    }
  } catch (error) {
    hide();
  }
  return false;
};

export const handleSyncCloud = async (fields: number | undefined) => {
  const hide = message.loading('同步中......', 30);
  if (fields === undefined) {
    hide();
    message.error("数据错误:不存在的数据记录")
    return false;
  }
  try {
    const {success} = await syncProvider(fields);
    hide();
    if (success) {
      message.success('同步成功');
      return true;
    }
  } catch (error) {
    hide();
  }
  return false;
}
