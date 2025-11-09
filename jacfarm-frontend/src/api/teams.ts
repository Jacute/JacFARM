import { API_URI, STATUS_OK, type BaseResponse, PAGE_LIMIT } from "./api";
import type { GetShortTeamsResponse, GetTeamsResponse } from "./responses";

export const getShortTeams = async () => {
    const res = await fetch(`${API_URI}/api/v1/teams/short`, {credentials: "include"});
    const data: GetShortTeamsResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return data.teams;
}

export const addTeam = async (name: string, ip: string) => {
    const res = await fetch(`${API_URI}/api/v1/teams`, {
        credentials: "include",
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: name,
            ip: ip
        })
    })
    const data: BaseResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
}

export const getTeams = async (page: number, query: string) => {
    const res = await fetch(`${API_URI}/api/v1/teams?limit=${PAGE_LIMIT}&page=${page}&query=${query}`, {credentials: "include"});
    const data: GetTeamsResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return {teams: data.teams, count: data.count};
}

export const deleteTeam = async (id: number) => {
    const res = await fetch(`${API_URI}/api/v1/teams/${id}`, {
        method: "DELETE",
        credentials: "include"
    });
    const data: BaseResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }
}