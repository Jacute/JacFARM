import type { Config } from "../models/models";
import { API_URI, PAGE_LIMIT, STATUS_OK, type BaseResponse } from "./api";
import type { GetConfigResponse } from "./responses";


export const getConfig = async (page: number): Promise<[Config[], number]> => {
    const res = await fetch(`${API_URI}/api/v1/config?page=${page}&limit=${PAGE_LIMIT}`, {credentials: "include"});
    const data: GetConfigResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return [data.config, data.count];
}

export const updateConfig = async (id: number, value: string) => {
    const res = await fetch(`${API_URI}/api/v1/config/${id}`, {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({value: value})
    });
    const data: BaseResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
}