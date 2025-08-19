import { API_URI, STATUS_OK, PAGE_LIMIT } from "./api";
import type { GetFlagsResponse, GetStatusesResponse } from "./responses";

export const getFlags = async (page: number, team_id: number, exploit_id: string, status_id: number) => {
    const res = await fetch(`${API_URI}/api/v1/flags?page=${page}&limit=${PAGE_LIMIT}&team_id=${team_id}&exploit_id=${exploit_id}&status_id=${status_id}`, {credentials: "include"});
    const data: GetFlagsResponse = await res.json();
    
    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return {flags: data.flags, count: data.count};
}

export const getFlagStatuses = async () => {
    const res = await fetch(`${API_URI}/api/v1/flags/statuses`, {credentials: "include"});
    const data: GetStatusesResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return data.statuses;
}

export const sendFlag = async (flag: string) => {
    const res = await fetch(`${API_URI}/api/v1/flags`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({flag: flag})
    });
    const data = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
}