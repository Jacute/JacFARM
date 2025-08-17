import { API_URI, STATUS_OK } from "./api";
import type { GetFlagsResponse } from "./responses";

export const pageLimit = 20;

export const getFlags = async (page: number) => {
    const res = await fetch(`${API_URI}/api/v1/flags?page=${page}&limit=${pageLimit}`, {credentials: "include"});
    const data: GetFlagsResponse = await res.json();
    
    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return data.flags;
}