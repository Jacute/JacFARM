import { API_URI, STATUS_OK, PAGE_LIMIT } from "./api";
import type { GetFlagsResponse, GetStatusesResponse } from "./responses";

export const getFlags = async (page: number, team_id: number | null, exploit_id: string | null, status_id: number | null) => {
    if (team_id == null) team_id = 0;
    if (exploit_id == null) exploit_id = "";
    if (status_id == null) status_id = 0;

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

export const getFlagsCount = async () => {
    const res = await fetch(`${API_URI}/api/v1/flags/count`, {credentials: "include"});
    const data = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return data.count;
}