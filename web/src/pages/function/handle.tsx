import {message} from "antd";
import type {FormValueType} from "./components/CreateForm";
import {createFunctionTunnel, deleteFunctionTunnel, updateFunctionTunnel} from "@/services/function/api";

export const handleCreateTunnel = async (fields: FormValueType) => {
  const hide = message.loading('创建中......', 30);
  try {
    const {success} = await createFunctionTunnel(fields);
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

export const handleUpdateTunnel = async (fields: FormValueType) => {
  const hide = message.loading('更新中......');
  try {
    const {success} = await updateFunctionTunnel(fields);
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

export const handleDeleteTunnel = async (fields: number | undefined) => {
  const hide = message.loading('删除中......');
  if (fields === undefined) {
    hide();
    message.error("数据错误:不存在的数据记录")
    return false;
  }
  try {
    const {success} = await deleteFunctionTunnel(fields);
    hide();
    if (success) {
      message.success('删除成功');
      return true;
    }
  } catch (error) {
    hide();
  }
  return false;
}
