import { API_URI, STATUS_OK } from "./api";
import type { GetShortTeamsResponse } from "./responses";

export const getShortTeams = async () => {
    const res = await fetch(`${API_URI}/api/v1/teams/short`, {credentials: "include"});
    const data: GetShortTeamsResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return data.teams;
}