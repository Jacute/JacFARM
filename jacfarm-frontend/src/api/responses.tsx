import type { Config, Exploit, ExploitShort, Flag, Status, Team, TeamShort } from "../models/models";
import type { BaseResponse } from "./api";

export interface GetFlagsResponse extends BaseResponse {
    flags: Array<Flag>
    count: number
}

export interface GetStatusesResponse extends BaseResponse {
    statuses: Array<Status>
}

export interface GetExploitsShortResponse extends BaseResponse {
    exploits: Array<ExploitShort>
}

export interface GetExploitsResponse extends BaseResponse {
    exploits: Array<Exploit>
    count: number
}

export interface GetShortTeamsResponse extends BaseResponse {
    teams: Array<TeamShort>
}

export interface GetTeamsResponse extends BaseResponse {
    teams: Array<Team>
    count: number
}

export interface GetConfigResponse extends BaseResponse {
    config: Config[]
    count: number
}