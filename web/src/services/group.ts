import { IGroup } from "../interfaces/group";
import api from "./api";

export const getMyGroups = async () => {
  return api.get("/user/me/group").then((res) => (res.data || []) as IGroup[]);
};

export const CreateGroup = async (payload: {name: string}) => {
  return api.post("/group", payload).then((res) => res.data as IGroup);
}
